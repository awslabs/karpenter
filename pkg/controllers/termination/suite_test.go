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

package termination_test

import (
	"fmt"
	"testing"

	"github.com/Pallinder/go-randomdata"
	"github.com/awslabs/karpenter/pkg/apis/provisioning/v1alpha3"
	"github.com/awslabs/karpenter/pkg/cloudprovider/fake"
	"github.com/awslabs/karpenter/pkg/cloudprovider/registry"
	"github.com/awslabs/karpenter/pkg/controllers/termination"
	"github.com/awslabs/karpenter/pkg/test"
	"github.com/awslabs/karpenter/pkg/utils/context"
	"go.uber.org/zap/zapcore"
	"sigs.k8s.io/controller-runtime/pkg/client"

	. "github.com/awslabs/karpenter/pkg/test/expectations"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	v1 "k8s.io/api/core/v1"
	corev1 "k8s.io/client-go/kubernetes/typed/core/v1"
	"knative.dev/pkg/ptr"
)

func TestAPIs(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Termination")
}

var ctx = context.NewLoggingContext(zapcore.DebugLevel)
var controller *termination.Controller
var evictionQueue *termination.EvictionQueue

var env = test.NewEnvironment(ctx, func(e *test.Environment) {
	cloudProvider := &fake.CloudProvider{}
	registry.RegisterOrDie(cloudProvider)
	coreV1Client := corev1.NewForConfigOrDie(e.Config)
	evictionQueue = termination.NewEvictionQueue(coreV1Client)
	controller = &termination.Controller{
		KubeClient: e.Client,
		Terminator: &termination.Terminator{
			KubeClient:    e.Client,
			CoreV1Client:  coreV1Client,
			CloudProvider: cloudProvider,
			EvictionQueue: evictionQueue,
		},
	}
})

var _ = BeforeSuite(func() {
	Expect(env.Start()).To(Succeed(), "Failed to start environment")
})

var _ = AfterSuite(func() {
	Expect(env.Stop()).To(Succeed(), "Failed to stop environment")
})

var _ = Describe("Termination", func() {
	var node *v1.Node

	BeforeEach(func() {
		node = test.Node(test.NodeOptions{Finalizers: []string{v1alpha3.KarpenterFinalizer}})
	})

	AfterEach(func() {
		ExpectCleanedUp(env.Client)
	})

	Context("Reconciliation", func() {
		It("should terminate deleted nodes", func() {
			ExpectCreated(env.Client, node)
			Expect(env.Client.Delete(ctx, node)).To(Succeed())
			node = ExpectNodeExists(env.Client, node.Name)
			ExpectReconcileSucceeded(ctx, controller, client.ObjectKeyFromObject(node))
			ExpectNotFound(env.Client, node)
		})
		It("should not evict pods that tolerate unschedulable taint", func() {
			podEvict := test.Pod(test.PodOptions{NodeName: node.Name})
			podSkip := test.Pod(test.PodOptions{
				NodeName:    node.Name,
				Tolerations: []v1.Toleration{{Key: v1.TaintNodeUnschedulable, Operator: v1.TolerationOpExists, Effect: v1.TaintEffectNoSchedule}},
			})
			ExpectCreated(env.Client, node, podEvict, podSkip)

			// Trigger Termination Controller
			Expect(env.Client.Delete(ctx, node)).To(Succeed())
			node = ExpectNodeExists(env.Client, node.Name)
			ExpectReconcileSucceeded(ctx, controller, client.ObjectKeyFromObject(node))

			// Expect podEvict to be enqueued for eviction
			ExpectEvicting(evictionQueue, podEvict)

			// Expect node to exist, but be cordoned
			node = ExpectNodeExists(env.Client, node.Name)
			Expect(node.Spec.Unschedulable).To(BeTrue())

			// Expect podEvict to be evicting, and delete it
			ExpectEvictingSucceeded(env.Client, podEvict)
			ExpectDeleted(env.Client, podEvict)

			// Reconcile to delete node
			node = ExpectNodeExists(env.Client, node.Name)
			ExpectReconcileSucceeded(ctx, controller, client.ObjectKeyFromObject(node))
			ExpectNotFound(env.Client, node)
		})
		It("should not terminate nodes that have a do-not-evict pod", func() {
			podEvict := test.Pod(test.PodOptions{NodeName: node.Name})
			podNoEvict := test.Pod(test.PodOptions{
				NodeName:    node.Name,
				Annotations: map[string]string{v1alpha3.KarpenterDoNotEvictPodAnnotation: "true"},
			})

			ExpectCreated(env.Client, node, podEvict, podNoEvict)

			Expect(env.Client.Delete(ctx, node)).To(Succeed())
			node = ExpectNodeExists(env.Client, node.Name)
			ExpectReconcileSucceeded(ctx, controller, client.ObjectKeyFromObject(node))

			// Expect no pod to be enqueued for eviction
			ExpectNotEvicting(evictionQueue, podEvict, podNoEvict)

			// Expect node to exist, but be cordoned
			node = ExpectNodeExists(env.Client, node.Name)
			Expect(node.Spec.Unschedulable).To(BeTrue())

			// Delete do-not-evict pod
			ExpectDeleted(env.Client, podNoEvict)

			// Reconcile node to evict pod
			node = ExpectNodeExists(env.Client, node.Name)
			ExpectReconcileSucceeded(ctx, controller, client.ObjectKeyFromObject(node))

			// Expect podEvict to be enqueued for eviction then be successful
			ExpectEvicting(evictionQueue, podEvict)
			ExpectEvictingSucceeded(env.Client, podEvict)

			// Delete pod to simulate successful eviction
			ExpectDeleted(env.Client, podEvict)

			// Reconcile to delete node
			node = ExpectNodeExists(env.Client, node.Name)
			ExpectReconcileSucceeded(ctx, controller, client.ObjectKeyFromObject(node))
			ExpectNotFound(env.Client, node)
		})
		It("should fail to evict pods that violate a PDB", func() {
			key, value := randomdata.SillyName(), randomdata.SillyName()
			pdb := test.PodDisruptionBudget(test.PDBOptions{
				Labels: map[string]string{key: value},
				// Don't let any pod evict
				MinAvailableNum: ptr.Int64(1),
			})
			podNoEvict := test.Pod(test.PodOptions{
				NodeName: node.Name,
				Labels:   map[string]string{key: value},
			})

			ExpectCreated(env.Client, node, podNoEvict, pdb)

			// Trigger Termination Controller
			Expect(env.Client.Delete(ctx, node)).To(Succeed())
			node = ExpectNodeExists(env.Client, node.Name)
			ExpectReconcileSucceeded(ctx, controller, client.ObjectKeyFromObject(node))

			// Expect the pod to be enqueued for eviction
			ExpectEvicting(evictionQueue, podNoEvict)

			// Expect node to exist, but be cordoned
			node = ExpectNodeExists(env.Client, node.Name)
			Expect(node.Spec.Unschedulable).To(BeTrue())

			// Expect podNoEvict to fail eviction due to PDB
			ExpectEvictingFailed(env.Client, evictionQueue, podNoEvict)

			// Delete pod to simulate successful eviction
			ExpectDeleted(env.Client, podNoEvict)

			// Reconcile to delete node
			node = ExpectNodeExists(env.Client, node.Name)
			ExpectReconcileSucceeded(ctx, controller, client.ObjectKeyFromObject(node))
			ExpectNotFound(env.Client, node)
		})
	})
})

func ExpectEvicting(e *termination.EvictionQueue, pods ...*v1.Pod) {
	for _, pod := range pods {
		Expect(e.Contains(client.ObjectKeyFromObject(pod))).To(BeTrue())
	}
}

func ExpectNotEvicting(e *termination.EvictionQueue, pods ...*v1.Pod) {
	for _, pod := range pods {
		Expect(e.Contains(client.ObjectKeyFromObject(pod))).To(BeFalse())
	}
}

func ExpectEvictingSucceeded(c client.Client, pods ...*v1.Pod) {
	for _, pod := range pods {
		Eventually(func() bool {
			return ExpectPodExists(c, pod.Name, pod.Namespace).GetDeletionTimestamp().IsZero()
		}, ReconcilerPropagationTime, RequestInterval).Should(BeFalse(), func() string {
			return fmt.Sprintf("expected %s/%s to be evicting, but it isn't", pod.Namespace, pod.Name)
		})
	}
}

func ExpectEvictingFailed(c client.Client, e *termination.EvictionQueue, pods ...*v1.Pod) {
	for _, pod := range pods {
		Eventually(func() bool {
			return ExpectPodExists(c, pod.Name, pod.Namespace).GetDeletionTimestamp().IsZero() && e.NumRequeues(client.ObjectKeyFromObject(pod)) > 0
		}, ReconcilerPropagationTime, RequestInterval).Should(BeTrue(), func() string {
			return fmt.Sprintf("expected %s/%s to not be evicting, but it is", pod.Namespace, pod.Name)
		})
	}
}
