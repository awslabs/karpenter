/*
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at
    http://www.apache.org/licenses/LICENSE-2.0
Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package aws

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/ec2metadata"
	"github.com/aws/aws-sdk-go/aws/endpoints"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/ssm"
	"github.com/awslabs/karpenter/pkg/apis/provisioning/v1alpha3"
	"github.com/awslabs/karpenter/pkg/cloudprovider"
	"github.com/awslabs/karpenter/pkg/cloudprovider/aws/utils"
	"github.com/awslabs/karpenter/pkg/utils/parallel"
	"github.com/awslabs/karpenter/pkg/utils/project"
	v1 "k8s.io/api/core/v1"
	"knative.dev/pkg/apis"
	"knative.dev/pkg/logging"
)

const (
	// CreationQPS limits the number of requests per second to CreateFleet
	// https://docs.aws.amazon.com/AWSEC2/latest/APIReference/throttling.html#throttling-limits
	CreationQPS = 2
	// CreationBurst limits the additional burst requests.
	// https://docs.aws.amazon.com/AWSEC2/latest/APIReference/throttling.html#throttling-limits
	CreationBurst = 100
	// CacheTTL restricts QPS to AWS APIs to this interval for verifying setup
	// resources. This value represents the maximum eventual consistency between
	// AWS actual state and the controller's ability to provision those
	// resources. Cache hits enable faster provisioning and reduced API load on
	// AWS APIs, which can have a serious import on performance and scalability.
	// DO NOT CHANGE THIS VALUE WITHOUT DUE CONSIDERATION
	CacheTTL = 60 * time.Second
	// CacheCleanupInterval triggers cache cleanup (lazy eviction) at this interval.
	CacheCleanupInterval = 10 * time.Minute
	// ClusterTagKeyFormat is set on all Kubernetes owned resources.
	ClusterTagKeyFormat = "kubernetes.io/cluster/%s"
	// KarpenterTagKeyFormat is set on all Karpenter owned resources.
	KarpenterTagKeyFormat = "karpenter.sh/cluster/%s"
)

var (
	SupportedOperatingSystems = []string{
		v1alpha3.OperatingSystemLinux,
	}
	SupportedArchitectures = []string{
		v1alpha3.ArchitectureAmd64,
		v1alpha3.ArchitectureArm64,
	}
)

type CloudProvider struct {
	nodeAPI                *NodeFactory
	launchTemplateProvider *LaunchTemplateProvider
	subnetProvider         *SubnetProvider
	instanceTypeProvider   *InstanceTypeProvider
	instanceProvider       *InstanceProvider
	creationQueue          *parallel.WorkQueue
}

func NewCloudProvider(options cloudprovider.Options) *CloudProvider {
	sess := withUserAgent(session.Must(
		session.NewSession(request.WithRetryer(
			&aws.Config{STSRegionalEndpoint: endpoints.RegionalSTSEndpoint},
			utils.NewRetryer()))))
	if *sess.Config.Region == "" {
		logging.FromContext(context.Background()).Info("AWS region not configured, asking EC2 Instance Metadata Service")
		*sess.Config.Region = getRegionFromIMDS(sess)
	}
	logging.FromContext(context.Background()).Infof("Using AWS region %s", *sess.Config.Region)
	ec2api := ec2.New(sess)
	return &CloudProvider{
		nodeAPI: &NodeFactory{ec2api: ec2api},
		launchTemplateProvider: NewLaunchTemplateProvider(
			ec2api,
			NewAMIProvider(ssm.New(sess), options.ClientSet),
			NewSecurityGroupProvider(ec2api),
		),
		subnetProvider:       NewSubnetProvider(ec2api),
		instanceTypeProvider: NewInstanceTypeProvider(ec2api),
		instanceProvider:     &InstanceProvider{ec2api: ec2api},
		creationQueue:        parallel.NewWorkQueue(CreationQPS, CreationBurst),
	}
}

// get the current region from EC2 IMDS
func getRegionFromIMDS(sess *session.Session) string {
	region, err := ec2metadata.New(sess).Region()
	if err != nil {
		panic(fmt.Sprintf("Failed to call the metadata server's region API, %s", err.Error()))
	}
	return region
}

// withUserAgent adds a karpenter specific user-agent string to AWS session
func withUserAgent(sess *session.Session) *session.Session {
	userAgent := fmt.Sprintf("karpenter.sh-%s", project.Version)
	sess.Handlers.Build.PushBack(request.MakeAddToUserAgentFreeFormHandler(userAgent))
	return sess
}

// Create a node given the constraints.
func (c *CloudProvider) Create(ctx context.Context, provisioner *v1alpha3.Provisioner, packing *cloudprovider.Packing, callback func(*v1.Node) error) chan error {
	return c.creationQueue.Add(func() error {
		return c.create(ctx, provisioner, packing, callback)
	})
}

func (c *CloudProvider) create(ctx context.Context, provisioner *v1alpha3.Provisioner, packing *cloudprovider.Packing, callback func(*v1.Node) error) error {
	constraints := Constraints{*packing.Constraints}
	// 1. Get Subnets and constrain by zones
	subnets, err := c.subnetProvider.Get(ctx, provisioner, &constraints)
	if err != nil {
		return fmt.Errorf("getting zonal subnets, %w", err)
	}
	// 2. Get Launch Template
	launchTemplate, err := c.launchTemplateProvider.Get(ctx, provisioner, &constraints)
	if err != nil {
		return fmt.Errorf("getting launch template, %w", err)
	}
	// 3. Create instance
	instanceID, err := c.instanceProvider.Create(ctx, launchTemplate, packing.InstanceTypeOptions, subnets, constraints.GetCapacityType())
	if err != nil {
		return fmt.Errorf("launching instances, %w", err)
	}
	// 4. Convert to node
	node, err := c.nodeAPI.For(ctx, instanceID)
	if err != nil {
		return fmt.Errorf("constructing node, %w", err)
	}
	return callback(node)
}

func (c *CloudProvider) GetInstanceTypes(ctx context.Context) ([]cloudprovider.InstanceType, error) {
	return c.instanceTypeProvider.Get(ctx)
}

func (c *CloudProvider) Terminate(ctx context.Context, node *v1.Node) error {
	return c.instanceProvider.Terminate(ctx, node)
}

// Validate cloud provider specific components of the cluster spec
func (c *CloudProvider) ValidateConstraints(ctx context.Context, constraints *v1alpha3.Constraints) (errs *apis.FieldError) {
	awsConstraints := Constraints{*constraints}
	return awsConstraints.Validate(ctx)
}

// Validate cloud provider specific components of the cluster spec.
func (c *CloudProvider) ValidateSpec(ctx context.Context, spec *v1alpha3.ProvisionerSpec) (errs *apis.FieldError) {
	if spec.Cluster.Name == nil || len(*spec.Cluster.Name) == 0 {
		errs = errs.Also(apis.ErrMissingField("name")).ViaField("cluster")
	}
	return errs
}
