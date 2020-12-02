// +build !ignore_autogenerated

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

// Code generated by controller-gen. DO NOT EDIT.

package v1alpha1

import (
	"k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"knative.dev/pkg/apis"
)

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Behavior) DeepCopyInto(out *Behavior) {
	*out = *in
	if in.ScaleUp != nil {
		in, out := &in.ScaleUp, &out.ScaleUp
		*out = new(ScalingRules)
		(*in).DeepCopyInto(*out)
	}
	if in.ScaleDown != nil {
		in, out := &in.ScaleDown, &out.ScaleDown
		*out = new(ScalingRules)
		(*in).DeepCopyInto(*out)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Behavior.
func (in *Behavior) DeepCopy() *Behavior {
	if in == nil {
		return nil
	}
	out := new(Behavior)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *CrossVersionObjectReference) DeepCopyInto(out *CrossVersionObjectReference) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new CrossVersionObjectReference.
func (in *CrossVersionObjectReference) DeepCopy() *CrossVersionObjectReference {
	if in == nil {
		return nil
	}
	out := new(CrossVersionObjectReference)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *HorizontalAutoscaler) DeepCopyInto(out *HorizontalAutoscaler) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	in.Status.DeepCopyInto(&out.Status)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new HorizontalAutoscaler.
func (in *HorizontalAutoscaler) DeepCopy() *HorizontalAutoscaler {
	if in == nil {
		return nil
	}
	out := new(HorizontalAutoscaler)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *HorizontalAutoscaler) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *HorizontalAutoscalerList) DeepCopyInto(out *HorizontalAutoscalerList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]HorizontalAutoscaler, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new HorizontalAutoscalerList.
func (in *HorizontalAutoscalerList) DeepCopy() *HorizontalAutoscalerList {
	if in == nil {
		return nil
	}
	out := new(HorizontalAutoscalerList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *HorizontalAutoscalerList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *HorizontalAutoscalerSpec) DeepCopyInto(out *HorizontalAutoscalerSpec) {
	*out = *in
	out.ScaleTargetRef = in.ScaleTargetRef
	if in.Metrics != nil {
		in, out := &in.Metrics, &out.Metrics
		*out = make([]Metric, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	in.Behavior.DeepCopyInto(&out.Behavior)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new HorizontalAutoscalerSpec.
func (in *HorizontalAutoscalerSpec) DeepCopy() *HorizontalAutoscalerSpec {
	if in == nil {
		return nil
	}
	out := new(HorizontalAutoscalerSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *HorizontalAutoscalerStatus) DeepCopyInto(out *HorizontalAutoscalerStatus) {
	*out = *in
	if in.LastScaleTime != nil {
		in, out := &in.LastScaleTime, &out.LastScaleTime
		*out = new(apis.VolatileTime)
		(*in).DeepCopyInto(*out)
	}
	if in.CurrentReplicas != nil {
		in, out := &in.CurrentReplicas, &out.CurrentReplicas
		*out = new(int32)
		**out = **in
	}
	if in.DesiredReplicas != nil {
		in, out := &in.DesiredReplicas, &out.DesiredReplicas
		*out = new(int32)
		**out = **in
	}
	if in.CurrentMetrics != nil {
		in, out := &in.CurrentMetrics, &out.CurrentMetrics
		*out = make([]MetricStatus, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	if in.Conditions != nil {
		in, out := &in.Conditions, &out.Conditions
		*out = make(apis.Conditions, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new HorizontalAutoscalerStatus.
func (in *HorizontalAutoscalerStatus) DeepCopy() *HorizontalAutoscalerStatus {
	if in == nil {
		return nil
	}
	out := new(HorizontalAutoscalerStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Metric) DeepCopyInto(out *Metric) {
	*out = *in
	if in.Prometheus != nil {
		in, out := &in.Prometheus, &out.Prometheus
		*out = new(PrometheusMetricSource)
		(*in).DeepCopyInto(*out)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Metric.
func (in *Metric) DeepCopy() *Metric {
	if in == nil {
		return nil
	}
	out := new(Metric)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *MetricStatus) DeepCopyInto(out *MetricStatus) {
	*out = *in
	if in.Object != nil {
		in, out := &in.Object, &out.Object
		*out = new(PrometheusMetricStatus)
		(*in).DeepCopyInto(*out)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new MetricStatus.
func (in *MetricStatus) DeepCopy() *MetricStatus {
	if in == nil {
		return nil
	}
	out := new(MetricStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *MetricTarget) DeepCopyInto(out *MetricTarget) {
	*out = *in
	if in.Value != nil {
		in, out := &in.Value, &out.Value
		x := (*in).DeepCopy()
		*out = &x
	}
	if in.AverageValue != nil {
		in, out := &in.AverageValue, &out.AverageValue
		x := (*in).DeepCopy()
		*out = &x
	}
	if in.AverageUtilization != nil {
		in, out := &in.AverageUtilization, &out.AverageUtilization
		*out = new(int32)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new MetricTarget.
func (in *MetricTarget) DeepCopy() *MetricTarget {
	if in == nil {
		return nil
	}
	out := new(MetricTarget)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *MetricValueStatus) DeepCopyInto(out *MetricValueStatus) {
	*out = *in
	if in.Value != nil {
		in, out := &in.Value, &out.Value
		x := (*in).DeepCopy()
		*out = &x
	}
	if in.AverageValue != nil {
		in, out := &in.AverageValue, &out.AverageValue
		x := (*in).DeepCopy()
		*out = &x
	}
	if in.AverageUtilization != nil {
		in, out := &in.AverageUtilization, &out.AverageUtilization
		*out = new(int32)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new MetricValueStatus.
func (in *MetricValueStatus) DeepCopy() *MetricValueStatus {
	if in == nil {
		return nil
	}
	out := new(MetricValueStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *MetricsProducer) DeepCopyInto(out *MetricsProducer) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	in.Status.DeepCopyInto(&out.Status)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new MetricsProducer.
func (in *MetricsProducer) DeepCopy() *MetricsProducer {
	if in == nil {
		return nil
	}
	out := new(MetricsProducer)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *MetricsProducer) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *MetricsProducerList) DeepCopyInto(out *MetricsProducerList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]MetricsProducer, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new MetricsProducerList.
func (in *MetricsProducerList) DeepCopy() *MetricsProducerList {
	if in == nil {
		return nil
	}
	out := new(MetricsProducerList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *MetricsProducerList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *MetricsProducerSpec) DeepCopyInto(out *MetricsProducerSpec) {
	*out = *in
	if in.PendingCapacity != nil {
		in, out := &in.PendingCapacity, &out.PendingCapacity
		*out = new(PendingCapacitySpec)
		(*in).DeepCopyInto(*out)
	}
	if in.Queue != nil {
		in, out := &in.Queue, &out.Queue
		*out = new(QueueSpec)
		**out = **in
	}
	if in.ReservedCapacity != nil {
		in, out := &in.ReservedCapacity, &out.ReservedCapacity
		*out = new(ReservedCapacitySpec)
		(*in).DeepCopyInto(*out)
	}
	if in.ScheduledCapacity != nil {
		in, out := &in.ScheduledCapacity, &out.ScheduledCapacity
		*out = new(ScheduledCapacitySpec)
		(*in).DeepCopyInto(*out)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new MetricsProducerSpec.
func (in *MetricsProducerSpec) DeepCopy() *MetricsProducerSpec {
	if in == nil {
		return nil
	}
	out := new(MetricsProducerSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *MetricsProducerStatus) DeepCopyInto(out *MetricsProducerStatus) {
	*out = *in
	if in.PendingCapacity != nil {
		in, out := &in.PendingCapacity, &out.PendingCapacity
		*out = new(PendingCapacityStatus)
		**out = **in
	}
	if in.Queue != nil {
		in, out := &in.Queue, &out.Queue
		*out = new(QueueStatus)
		**out = **in
	}
	if in.ReservedCapacity != nil {
		in, out := &in.ReservedCapacity, &out.ReservedCapacity
		*out = make(map[v1.ResourceName]string, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
	if in.ScheduledCapacity != nil {
		in, out := &in.ScheduledCapacity, &out.ScheduledCapacity
		*out = new(ScheduledCapacityStatus)
		**out = **in
	}
	if in.Conditions != nil {
		in, out := &in.Conditions, &out.Conditions
		*out = make(apis.Conditions, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new MetricsProducerStatus.
func (in *MetricsProducerStatus) DeepCopy() *MetricsProducerStatus {
	if in == nil {
		return nil
	}
	out := new(MetricsProducerStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *PendingCapacitySpec) DeepCopyInto(out *PendingCapacitySpec) {
	*out = *in
	if in.NodeSelector != nil {
		in, out := &in.NodeSelector, &out.NodeSelector
		*out = make(map[string]string, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new PendingCapacitySpec.
func (in *PendingCapacitySpec) DeepCopy() *PendingCapacitySpec {
	if in == nil {
		return nil
	}
	out := new(PendingCapacitySpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *PendingCapacityStatus) DeepCopyInto(out *PendingCapacityStatus) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new PendingCapacityStatus.
func (in *PendingCapacityStatus) DeepCopy() *PendingCapacityStatus {
	if in == nil {
		return nil
	}
	out := new(PendingCapacityStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *PendingPodsSpec) DeepCopyInto(out *PendingPodsSpec) {
	*out = *in
	if in.NodeSelector != nil {
		in, out := &in.NodeSelector, &out.NodeSelector
		*out = make(map[string]string, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new PendingPodsSpec.
func (in *PendingPodsSpec) DeepCopy() *PendingPodsSpec {
	if in == nil {
		return nil
	}
	out := new(PendingPodsSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *PrometheusMetricSource) DeepCopyInto(out *PrometheusMetricSource) {
	*out = *in
	in.Target.DeepCopyInto(&out.Target)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new PrometheusMetricSource.
func (in *PrometheusMetricSource) DeepCopy() *PrometheusMetricSource {
	if in == nil {
		return nil
	}
	out := new(PrometheusMetricSource)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *PrometheusMetricStatus) DeepCopyInto(out *PrometheusMetricStatus) {
	*out = *in
	in.Current.DeepCopyInto(&out.Current)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new PrometheusMetricStatus.
func (in *PrometheusMetricStatus) DeepCopy() *PrometheusMetricStatus {
	if in == nil {
		return nil
	}
	out := new(PrometheusMetricStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *QueueSpec) DeepCopyInto(out *QueueSpec) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new QueueSpec.
func (in *QueueSpec) DeepCopy() *QueueSpec {
	if in == nil {
		return nil
	}
	out := new(QueueSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *QueueStatus) DeepCopyInto(out *QueueStatus) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new QueueStatus.
func (in *QueueStatus) DeepCopy() *QueueStatus {
	if in == nil {
		return nil
	}
	out := new(QueueStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ReservedCapacitySpec) DeepCopyInto(out *ReservedCapacitySpec) {
	*out = *in
	if in.NodeSelector != nil {
		in, out := &in.NodeSelector, &out.NodeSelector
		*out = make(map[string]string, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ReservedCapacitySpec.
func (in *ReservedCapacitySpec) DeepCopy() *ReservedCapacitySpec {
	if in == nil {
		return nil
	}
	out := new(ReservedCapacitySpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ScalableNodeGroup) DeepCopyInto(out *ScalableNodeGroup) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	in.Status.DeepCopyInto(&out.Status)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ScalableNodeGroup.
func (in *ScalableNodeGroup) DeepCopy() *ScalableNodeGroup {
	if in == nil {
		return nil
	}
	out := new(ScalableNodeGroup)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *ScalableNodeGroup) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ScalableNodeGroupList) DeepCopyInto(out *ScalableNodeGroupList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]ScalableNodeGroup, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ScalableNodeGroupList.
func (in *ScalableNodeGroupList) DeepCopy() *ScalableNodeGroupList {
	if in == nil {
		return nil
	}
	out := new(ScalableNodeGroupList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *ScalableNodeGroupList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ScalableNodeGroupSpec) DeepCopyInto(out *ScalableNodeGroupSpec) {
	*out = *in
	if in.Replicas != nil {
		in, out := &in.Replicas, &out.Replicas
		*out = new(int32)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ScalableNodeGroupSpec.
func (in *ScalableNodeGroupSpec) DeepCopy() *ScalableNodeGroupSpec {
	if in == nil {
		return nil
	}
	out := new(ScalableNodeGroupSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ScalableNodeGroupStatus) DeepCopyInto(out *ScalableNodeGroupStatus) {
	*out = *in
	if in.Replicas != nil {
		in, out := &in.Replicas, &out.Replicas
		*out = new(int32)
		**out = **in
	}
	if in.Conditions != nil {
		in, out := &in.Conditions, &out.Conditions
		*out = make(apis.Conditions, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ScalableNodeGroupStatus.
func (in *ScalableNodeGroupStatus) DeepCopy() *ScalableNodeGroupStatus {
	if in == nil {
		return nil
	}
	out := new(ScalableNodeGroupStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ScalingPolicy) DeepCopyInto(out *ScalingPolicy) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ScalingPolicy.
func (in *ScalingPolicy) DeepCopy() *ScalingPolicy {
	if in == nil {
		return nil
	}
	out := new(ScalingPolicy)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ScalingRules) DeepCopyInto(out *ScalingRules) {
	*out = *in
	if in.StabilizationWindowSeconds != nil {
		in, out := &in.StabilizationWindowSeconds, &out.StabilizationWindowSeconds
		*out = new(int32)
		**out = **in
	}
	if in.SelectPolicy != nil {
		in, out := &in.SelectPolicy, &out.SelectPolicy
		*out = new(ScalingPolicySelect)
		**out = **in
	}
	if in.Policies != nil {
		in, out := &in.Policies, &out.Policies
		*out = make([]ScalingPolicy, len(*in))
		copy(*out, *in)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ScalingRules.
func (in *ScalingRules) DeepCopy() *ScalingRules {
	if in == nil {
		return nil
	}
	out := new(ScalingRules)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ScheduledBehavior) DeepCopyInto(out *ScheduledBehavior) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ScheduledBehavior.
func (in *ScheduledBehavior) DeepCopy() *ScheduledBehavior {
	if in == nil {
		return nil
	}
	out := new(ScheduledBehavior)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ScheduledCapacitySpec) DeepCopyInto(out *ScheduledCapacitySpec) {
	*out = *in
	if in.NodeSelector != nil {
		in, out := &in.NodeSelector, &out.NodeSelector
		*out = make(map[string]string, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
	if in.Behaviors != nil {
		in, out := &in.Behaviors, &out.Behaviors
		*out = make([]ScheduledBehavior, len(*in))
		copy(*out, *in)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ScheduledCapacitySpec.
func (in *ScheduledCapacitySpec) DeepCopy() *ScheduledCapacitySpec {
	if in == nil {
		return nil
	}
	out := new(ScheduledCapacitySpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ScheduledCapacityStatus) DeepCopyInto(out *ScheduledCapacityStatus) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ScheduledCapacityStatus.
func (in *ScheduledCapacityStatus) DeepCopy() *ScheduledCapacityStatus {
	if in == nil {
		return nil
	}
	out := new(ScheduledCapacityStatus)
	in.DeepCopyInto(out)
	return out
}
