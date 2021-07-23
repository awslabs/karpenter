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

package reallocation_test

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/Pallinder/go-randomdata"
	"github.com/awslabs/karpenter/pkg/apis/provisioning/v1alpha3"
	"github.com/awslabs/karpenter/pkg/cloudprovider/fake"
	"github.com/awslabs/karpenter/pkg/cloudprovider/registry"
	"github.com/awslabs/karpenter/pkg/controllers/reallocation"
	"github.com/awslabs/karpenter/pkg/test"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	corev1 "k8s.io/client-go/kubernetes/typed/core/v1"
	"knative.dev/pkg/ptr"

	. "github.com/awslabs/karpenter/pkg/test/expectations"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func TestAPIs(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Provisioner/Reallocator")
}

var controller *reallocation.Controller
var env = test.NewEnvironment(func(e *test.Environment) {
	cloudProvider := &fake.CloudProvider{}
	registry.RegisterOrDie(cloudProvider)
	controller = reallocation.NewController(
		e.Client,
		corev1.NewForConfigOrDie(e.Config),
		cloudProvider,
	)
})

var _ = BeforeSuite(func() {
	Expect(env.Start()).To(Succeed(), "Failed to start environment")
})

var _ = AfterSuite(func() {
	Expect(env.Stop()).To(Succeed(), "Failed to stop environment")
})

var _ = Describe("Reallocation", func() {
	var provisioner *v1alpha3.Provisioner
	var ctx context.Context

	BeforeEach(func() {
		provisioner = &v1alpha3.Provisioner{
			ObjectMeta: metav1.ObjectMeta{Name: v1alpha3.DefaultProvisioner.Name},
			Spec: v1alpha3.ProvisionerSpec{
				Cluster:              v1alpha3.Cluster{Name: ptr.String("test-cluster"), Endpoint: "http://test-cluster", CABundle: ptr.String("dGVzdC1jbHVzdGVyCg==")},
				TTLSecondsAfterEmpty: ptr.Int64(300),
			},
		}
		ctx = context.Background()
	})

	AfterEach(func() {
		ExpectCleanedUp(env.Client)
	})

	Context("Reconciliation", func() {
		It("should not TTL nodes that have ready status unknown", func() {
			node := test.Node(test.NodeOptions{
				ReadyStatus: v1.ConditionUnknown,
			})

			ExpectCreated(env.Client, provisioner)
			ExpectCreatedWithStatus(env.Client, node)
			ExpectReconcileSucceeded(controller, client.ObjectKeyFromObject(provisioner))

			updatedNode := &v1.Node{}
			Expect(env.Client.Get(ctx, client.ObjectKey{Name: node.Name}, updatedNode)).To(Succeed())
			Expect(updatedNode.Labels).ToNot(HaveKey(v1alpha3.ProvisionerUnderutilizedLabelKey))
			Expect(updatedNode.Annotations).ToNot(HaveKey(v1alpha3.ProvisionerTTLAfterEmptyKey))
		})
		It("should not TTL nodes that have ready status unknown", func() {
			node := test.Node(test.NodeOptions{
				ReadyStatus: v1.ConditionFalse,
			})

			ExpectCreated(env.Client, provisioner)
			ExpectCreatedWithStatus(env.Client, node)
			ExpectReconcileSucceeded(controller, client.ObjectKeyFromObject(provisioner))

			updatedNode := &v1.Node{}
			Expect(env.Client.Get(ctx, client.ObjectKey{Name: node.Name}, updatedNode)).To(Succeed())
			Expect(updatedNode.Labels).ToNot(HaveKey(v1alpha3.ProvisionerUnderutilizedLabelKey))
			Expect(updatedNode.Annotations).ToNot(HaveKey(v1alpha3.ProvisionerTTLAfterEmptyKey))
		})
		It("should label nodes as underutilized and add TTL", func() {
			node := test.Node(test.NodeOptions{
				Labels: map[string]string{
					v1alpha3.ProvisionerNameLabelKey: provisioner.Name,
				},
			})
			ExpectCreated(env.Client, provisioner)
			ExpectCreatedWithStatus(env.Client, node)
			ExpectReconcileSucceeded(controller, client.ObjectKeyFromObject(provisioner))

			updatedNode := &v1.Node{}
			Expect(env.Client.Get(ctx, client.ObjectKey{Name: node.Name}, updatedNode)).To(Succeed())
			Expect(updatedNode.Labels).To(HaveKey(v1alpha3.ProvisionerUnderutilizedLabelKey))
			Expect(updatedNode.Annotations).To(HaveKey(v1alpha3.ProvisionerTTLAfterEmptyKey))
		})
		It("should remove labels from utilized nodes", func() {
			node := test.Node(test.NodeOptions{
				Labels: map[string]string{
					v1alpha3.ProvisionerNameLabelKey:          provisioner.Name,
					v1alpha3.ProvisionerUnderutilizedLabelKey: "true",
				},
				Annotations: map[string]string{
					v1alpha3.ProvisionerTTLAfterEmptyKey: time.Now().Add(time.Duration(100) * time.Second).Format(time.RFC3339),
				},
			})
			ExpectCreated(env.Client, provisioner)
			ExpectCreatedWithStatus(env.Client, node)
			ExpectCreatedWithStatus(env.Client, test.Pod(test.PodOptions{
				Name:       strings.ToLower(randomdata.SillyName()),
				Namespace:  provisioner.Namespace,
				NodeName:   node.Name,
				Conditions: []v1.PodCondition{{Type: v1.PodReady, Status: v1.ConditionTrue}},
			}))
			ExpectReconcileSucceeded(controller, client.ObjectKeyFromObject(provisioner))

			updatedNode := &v1.Node{}
			Expect(env.Client.Get(ctx, client.ObjectKey{Name: node.Name}, updatedNode)).To(Succeed())
			Expect(updatedNode.Labels).ToNot(HaveKey(v1alpha3.ProvisionerUnderutilizedLabelKey))
			Expect(updatedNode.Annotations).ToNot(HaveKey(v1alpha3.ProvisionerTTLAfterEmptyKey))
		})
	})
})
