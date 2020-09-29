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

package v1alpha1

import (
	"strings"
	"testing"

	"github.com/Pallinder/go-randomdata"
	v1alpha1 "github.com/ellistarn/karpenter/pkg/apis/autoscaling/v1alpha1"
	"github.com/ellistarn/karpenter/pkg/controllers"
	"github.com/ellistarn/karpenter/pkg/metrics/producers"
	"github.com/ellistarn/karpenter/pkg/test"
	"github.com/ellistarn/karpenter/pkg/test/environment"
	. "github.com/ellistarn/karpenter/pkg/test/expectations"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"sigs.k8s.io/controller-runtime/pkg/envtest/printer"
)

func TestAPIs(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecsWithDefaultAndCustomReporters(t,
		"Metrics Producer Suite",
		[]Reporter{printer.NewlineReporter{}})
}

func injectHorizontalAutoscalerController(environment *environment.Local) {
	controller := &Controller{
		ProducerFactory: producers.Factory{
			InformerFactory: environment.InformerFactory,
		},
		Client: environment.Manager.GetClient(),
	}
	Expect(controllers.RegisterController(environment.Manager, controller)).To(Succeed(), "Failed to register controller")
	Expect(controllers.RegisterWebhook(environment.Manager, controller)).To(Succeed(), "Failed to register webhook")
}

var env environment.Environment = environment.NewLocal(injectHorizontalAutoscalerController)

var _ = BeforeSuite(func() {
	Expect(env.Start()).To(Succeed(), "Failed to start environment")
})

var _ = AfterSuite(func() {
	Expect(env.Stop()).To(Succeed(), "Failed to stop environment")
})

var _ = Describe("Test Samples", func() {
	var ns *environment.Namespace
	var mp *v1alpha1.MetricsProducer

	BeforeEach(func() {
		var err error
		ns, err = env.NewNamespace()
		Expect(err).NotTo(HaveOccurred())
		mp = &v1alpha1.MetricsProducer{}
	})

	Context("Capacity Reservations", func() {
		It("should produce reservation metrics for 7/48 cores, 77/384 memory", func() {
			Expect(ns.ParseResources("docs/samples/reserved-capacity/resources.yaml", mp)).To(Succeed())
			mp.Spec.ReservedCapacity.NodeSelector = map[string]string{"k8s.io/nodegroup": strings.ToLower(randomdata.SillyName())}

			nodes := []test.Object{
				test.Node(mp.Spec.ReservedCapacity.NodeSelector, 16, 128),
				test.Node(mp.Spec.ReservedCapacity.NodeSelector, 16, 128),
				test.Node(mp.Spec.ReservedCapacity.NodeSelector, 16, 128),
			}
			pods := []test.Object{
				// node[0] 6/16 cores, 76/128 gig allocated
				test.Pod(nodes[0].GetName(), ns.Name, 1, 1),
				test.Pod(nodes[0].GetName(), ns.Name, 2, 25),
				test.Pod(nodes[0].GetName(), ns.Name, 3, 50),
				// node[1] 1/16 cores, 76/128 gig allocated
				test.Pod(nodes[1].GetName(), ns.Name, 1, 1),
				// node[2] is unallocated,
			}

			ExpectCreated(ns.Client, nodes...)
			ExpectCreated(ns.Client, pods...)

			ExpectEventuallyCreated(ns.Client, mp)
			// ExpectEventuallyHappy(ns.Client, mp)
			// TODO Verify metrics as set as expected
			ExpectEventuallyDeleted(ns.Client, mp)

			ExpectDeleted(ns.Client, nodes...)
			ExpectDeleted(ns.Client, pods...)
		})
	})
})
