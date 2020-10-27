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

package mock

import (
	"fmt"

	"github.com/ellistarn/karpenter/pkg/apis/autoscaling/v1alpha1"
	"github.com/ellistarn/karpenter/pkg/cloudprovider"
)

var (
	NotImplementedError = fmt.Errorf("provider is not implemented")
)

type Factory struct {
	WantErr error
}

func NewFactory() *Factory {
	return &Factory{}
}

func NewNotImplementedFactory() *Factory {
	return &Factory{WantErr: NotImplementedError}
}

func (f *Factory) NodeGroupFor(sng *v1alpha1.ScalableNodeGroupSpec) cloudprovider.NodeGroup {
	return &NodeGroup{Replicas: sng.Replicas, WantErr: f.WantErr}
}

func (f *Factory) QueueFor(spec *v1alpha1.QueueSpec) cloudprovider.Queue {
	return &Queue{Id: spec.ID, WantErr: f.WantErr}
}
