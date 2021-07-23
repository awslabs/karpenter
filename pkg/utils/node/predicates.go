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

package node

import (
	"time"

	"github.com/awslabs/karpenter/pkg/apis/provisioning/v1alpha3"
	"github.com/awslabs/karpenter/pkg/utils/pod"
	v1 "k8s.io/api/core/v1"
)

func IsReadyAndSchedulable(node v1.Node) bool {
	for _, condition := range node.Status.Conditions {
		if condition.Type == v1.NodeReady {
			return condition.Status == v1.ConditionTrue && !node.Spec.Unschedulable
		}
	}
	return false
}

func IsReady(node v1.Node) bool {
	for _, condition := range node.Status.Conditions {
		if condition.Type == v1.NodeReady {
			return condition.Status == v1.ConditionTrue
		}
	}
	return false
}

func FailsToBecomeReady(node v1.Node) bool {
	if time.Since(node.GetCreationTimestamp().Time) < (15 * time.Minute) {
		return false
	}
	for _, condition := range node.Status.Conditions {
		if condition.Type == v1.NodeReady {
			return condition.Status == v1.ConditionUnknown && condition.LastTransitionTime.IsZero() && condition.LastHeartbeatTime.IsZero()
		}
	}
	return false
}

func IsPastEmptyTTL(node v1.Node) bool {
	ttl, ok := node.Annotations[v1alpha3.ProvisionerTTLAfterEmptyKey]
	if !ok {
		return false
	}
	ttlTime, err := time.Parse(time.RFC3339, ttl)
	if err != nil {
		return false
	}
	return time.Now().After(ttlTime)
}

// IsEmpty returns if the node has 0 non-daemonset pods
func IsEmpty(pods []*v1.Pod) bool {
	for _, p := range pods {
		if pod.HasFailed(p) {
			continue
		}
		if !pod.IsOwnedByDaemonSet(p) {
			return false
		}
	}
	return true
}
