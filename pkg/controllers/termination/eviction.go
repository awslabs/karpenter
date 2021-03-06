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

package termination

import (
	"context"
	"sync"
	"time"

	set "github.com/deckarep/golang-set"
	"go.uber.org/zap"
	v1 "k8s.io/api/core/v1"
	"k8s.io/api/policy/v1beta1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	corev1 "k8s.io/client-go/kubernetes/typed/core/v1"
	"k8s.io/client-go/util/workqueue"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const (
	evictionQueueBaseDelay = 100 * time.Millisecond
	evictionQueueMaxDelay  = 10 * time.Second
)

type EvictionQueue struct {
	workqueue.RateLimitingInterface
	set.Set

	coreV1Client corev1.CoreV1Interface
	once         sync.Once
}

func NewEvictionQueue(coreV1Client corev1.CoreV1Interface) *EvictionQueue {
	return &EvictionQueue{
		RateLimitingInterface: workqueue.NewRateLimitingQueue(workqueue.NewItemExponentialFailureRateLimiter(evictionQueueBaseDelay, evictionQueueMaxDelay)),
		Set:                   set.NewSet(),

		coreV1Client: coreV1Client,
	}
}

// Add adds pods to the EvictionQueue
func (e *EvictionQueue) Add(pods []*v1.Pod) {
	// Start processing eviction queue if it hasn't started already
	e.once.Do(func() { go e.run() })

	for _, pod := range pods {
		if nn := client.ObjectKeyFromObject(pod); !e.Set.Contains(nn) {
			e.Set.Add(nn)
			e.RateLimitingInterface.Add(nn)
		}
	}
}

func (e *EvictionQueue) run() {
	for {
		// Get pod from queue. This waits until queue is non-empty.
		item, shutdown := e.RateLimitingInterface.Get()
		if shutdown {
			break
		}
		nn := item.(types.NamespacedName)
		// Evict pod
		if e.evict(nn) {
			zap.S().Debugf("Evicted pod %s", nn.String())
			e.RateLimitingInterface.Forget(nn)
			e.Set.Remove(nn)
			e.RateLimitingInterface.Done(nn)
			continue
		}
		e.RateLimitingInterface.Done(nn)
		// Requeue pod if eviction failed
		e.RateLimitingInterface.AddRateLimited(nn)
	}
	zap.S().Errorf("EvictionQueue is broken and has shutdown.")
}

// evict returns true if successful eviction call, error is returned if not eviction-related error
func (e *EvictionQueue) evict(nn types.NamespacedName) bool {
	err := e.coreV1Client.Pods(nn.Namespace).Evict(context.Background(), &v1beta1.Eviction{
		ObjectMeta: metav1.ObjectMeta{Name: nn.Name, Namespace: nn.Namespace},
	})
	if errors.IsInternalError(err) { // 500
		zap.S().Debugf("Failed to evict pod %s due to PDB misconfiguration error.", nn.String())
		return false
	}
	if errors.IsTooManyRequests(err) { // 429
		zap.S().Debugf("Failed to evict pod %s due to PDB violation.", nn.String())
		return false
	}
	if errors.IsNotFound(err) { // 404
		return true
	}
	if err != nil {
		return false
	}
	return true
}
