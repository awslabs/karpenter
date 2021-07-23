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

package v1alpha3

import (
	"strings"
	"testing"

	"github.com/Pallinder/go-randomdata"
	"github.com/awslabs/karpenter/pkg/utils/context"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"go.uber.org/zap/zapcore"
	"knative.dev/pkg/ptr"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var ctx = context.NewLoggingContext(zapcore.DebugLevel)

func TestAPIs(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Validation")
}

var _ = Describe("Validation", func() {
	var provisioner *Provisioner

	BeforeEach(func() {
		provisioner = &Provisioner{
			ObjectMeta: metav1.ObjectMeta{
				Name: strings.ToLower(randomdata.SillyName()),
			},
			Spec: ProvisionerSpec{
				Cluster: Cluster{
					Name:     ptr.String("test-cluster"),
					Endpoint: "https://test-cluster",
					CABundle: ptr.String("dGVzdC1jbHVzdGVyCg=="),
				},
			},
		}
	})

	It("should fail on negative expiry ttl", func() {
		provisioner.Spec.TTLSecondsUntilExpired = ptr.Int64(-1)
		Expect(provisioner.Validate(ctx)).ToNot(Succeed())
	})

	It("should fail on negative empty ttl", func() {
		provisioner.Spec.TTLSecondsAfterEmpty = ptr.Int64(-1)
		Expect(provisioner.Validate(ctx)).ToNot(Succeed())
	})

	It("should fail for empty cluster specification", func() {
		for _, cluster := range []Cluster{
			{},
			{Name: ptr.String("test-cluster"), CABundle: ptr.String("dGVzdC1jbHVzdGVyCg==")},
			{Name: ptr.String("test-cluster")},
			{CABundle: ptr.String("dGVzdC1jbHVzdGVyCg==")},
		} {
			provisioner.Spec.Cluster = cluster
			Expect(provisioner.Validate(ctx)).ToNot(Succeed())
		}
	})

	Context("Labels", func() {
		It("should fail for invalid label keys", func() {
			provisioner.Spec.Labels = map[string]string{"spaces are not allowed": randomdata.SillyName()}
			Expect(provisioner.Validate(ctx)).ToNot(Succeed())
		})
		It("should fail for invalid label values", func() {
			provisioner.Spec.Labels = map[string]string{randomdata.SillyName(): "/ is not allowed"}
			Expect(provisioner.Validate(ctx)).ToNot(Succeed())
		})
		It("should fail for restricted labels", func() {
			for _, label := range []string{
				ArchitectureLabelKey,
				OperatingSystemLabelKey,
				ProvisionerNameLabelKey,
				ProvisionerUnderutilizedLabelKey,
				ZoneLabelKey,
				InstanceTypeLabelKey,
			} {
				provisioner.Spec.Labels = map[string]string{label: randomdata.SillyName()}
				Expect(provisioner.Validate(ctx)).ToNot(Succeed())
			}
		})
	})
	Context("Taints", func() {
		It("should succeed for valid taints", func() {
			provisioner.Spec.Taints = []v1.Taint{
				{Key: "a", Value: "b", Effect: v1.TaintEffectNoSchedule},
				{Key: "c", Value: "d", Effect: v1.TaintEffectNoExecute},
				{Key: "e", Value: "f", Effect: v1.TaintEffectPreferNoSchedule},
				{Key: "key-only", Effect: v1.TaintEffectNoExecute},
			}
			Expect(provisioner.Validate(ctx)).To(Succeed())
		})
		It("should fail for invalid taint keys", func() {
			provisioner.Spec.Taints = []v1.Taint{{Key: "???"}}
			Expect(provisioner.Validate(ctx)).ToNot(Succeed())
		})
		It("should fail for missing taint key", func() {
			provisioner.Spec.Taints = []v1.Taint{{Effect: v1.TaintEffectNoSchedule}}
			Expect(provisioner.Validate(ctx)).ToNot(Succeed())
		})
		It("should fail for invalid taint value", func() {
			provisioner.Spec.Taints = []v1.Taint{{Key: "invalid-value", Effect: v1.TaintEffectNoSchedule, Value: "???"}}
			Expect(provisioner.Validate(ctx)).ToNot(Succeed())
		})
		It("should fail for invalid taint effect", func() {
			provisioner.Spec.Taints = []v1.Taint{{Key: "invalid-effect", Effect: "???"}}
			Expect(provisioner.Validate(ctx)).ToNot(Succeed())
		})
	})
	Context("Zones", func() {
		SupportedZones = append(SupportedZones, "test-zone-1")
		It("should succeed if unspecified", func() {
			Expect(provisioner.Validate(ctx)).To(Succeed())
		})
		It("should fail if not supported", func() {
			provisioner.Spec.Zones = []string{"unknown"}
			Expect(provisioner.Validate(ctx)).ToNot(Succeed())
		})
		It("should succeed if supported", func() {
			provisioner.Spec.Zones = []string{"test-zone-1"}
			Expect(provisioner.Validate(ctx)).To(Succeed())
		})
	})

	Context("InstanceTypes", func() {
		SupportedInstanceTypes = append(SupportedInstanceTypes, "test-instance-type")
		It("should succeed if unspecified", func() {
			Expect(provisioner.Validate(ctx)).To(Succeed())
		})
		It("should fail if not supported", func() {
			provisioner.Spec.InstanceTypes = []string{"unknown"}
			Expect(provisioner.Validate(ctx)).ToNot(Succeed())
		})
		It("should succeed if supported", func() {
			provisioner.Spec.InstanceTypes = []string{
				"test-instance-type",
			}
			Expect(provisioner.Validate(ctx)).To(Succeed())
		})
	})

	Context("Architecture", func() {
		SupportedArchitectures = append(SupportedArchitectures, "test-architecture")
		It("should succeed if unspecified", func() {
			Expect(provisioner.Validate(ctx)).To(Succeed())
		})
		It("should fail if not supported", func() {
			provisioner.Spec.Architecture = ptr.String("unknown")
			Expect(provisioner.Validate(ctx)).ToNot(Succeed())
		})
		It("should succeed if supported", func() {
			provisioner.Spec.Architecture = ptr.String("test-architecture")
			Expect(provisioner.Validate(ctx)).To(Succeed())
		})
	})

	Context("OperatingSystem", func() {
		SupportedOperatingSystems = append(SupportedArchitectures, "test-operating-system")
		It("should succeed if unspecified", func() {
			Expect(provisioner.Validate(ctx)).To(Succeed())
		})
		It("should fail if not supported", func() {
			provisioner.Spec.OperatingSystem = ptr.String("unknown")
			Expect(provisioner.Validate(ctx)).ToNot(Succeed())
		})
		It("should succeed if supported", func() {
			provisioner.Spec.OperatingSystem = ptr.String("test-operating-system")
			Expect(provisioner.Validate(ctx)).To(Succeed())
		})
	})
})
