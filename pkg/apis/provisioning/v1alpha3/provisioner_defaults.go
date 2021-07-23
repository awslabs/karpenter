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
	"context"
	"encoding/base64"
	"errors"
	"io/ioutil"
	"os"

	"knative.dev/pkg/logging"
)

const (
	InClusterCABundlePath = "/var/run/secrets/kubernetes.io/serviceaccount/ca.crt"
)

// SetDefaults for the provisioner
func (p *Provisioner) SetDefaults(ctx context.Context) {}

// SetDefaults for the provisioner, cascading to all subspecs
func (s *ProvisionerSpec) SetDefaults(ctx context.Context) {}

// WithDefaults returns a copy of this Provisioner with some empty/missing
// properties replaced by (potentially dynamic) cloud provider agnostic default values.
// The returned copy might be complemented by dynamic default values which
// must not be hoisted (saved) into the original Provisioner CRD as those
// default values might change over time (e.g. rolling upgrade of CABundle, ...).
func (p *Provisioner) WithDynamicDefaults(ctx context.Context) (_ Provisioner, err error) {
	provisioner := *p.DeepCopy()
	provisioner.Spec, err = provisioner.Spec.withDynamicDefaults(ctx)
	return provisioner, err
}

// WithDefaults returns a copy of this ProvisionerSpec with some empty
// properties replaced by default values.
func (s *ProvisionerSpec) withDynamicDefaults(ctx context.Context) (_ ProvisionerSpec, err error) {
	spec := *s.DeepCopy()
	spec.Cluster, err = spec.Cluster.withDynamicDefaults(ctx)
	return spec, err
}

// WithDefaults returns a copy of this Cluster with some empty
// properties replaced by default values. Notably, it will try
// to load the CABundle from the in-cluster configuraiton if it
// is not explicitly set.
func (c *Cluster) withDynamicDefaults(ctx context.Context) (_ Cluster, err error) {
	cluster := *c.DeepCopy()
	cluster.CABundle, err = cluster.getCABundle(ctx)
	return cluster, err
}

func (c *Cluster) getCABundle(ctx context.Context) (*string, error) {
	if c.CABundle != nil {
		// If CABundle is explicitly provided use that one. An empty string is
		// a valid value here if the intention is to disable the in-cluster CABundle
		// and using the HTTP client's default trust-store (CABundle) instead.
		logging.FromContext(ctx).Debugf("Using inline CABundle from Provisioner specification")
		return c.CABundle, nil
	}
	// Otherwise, fallback to the in-cluster configuration.
	binary, err := ioutil.ReadFile(InClusterCABundlePath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			logging.FromContext(ctx).Debugf("In-cluster CABundle file %s not found, will use HTTP client's default trust-store instead, %w", InClusterCABundlePath, err)
			return nil, nil
		}
		return nil, err
	}
	logging.FromContext(ctx).Debugf("Using in-cluster CABundle from file %s", InClusterCABundlePath)
	encoded := base64.StdEncoding.EncodeToString(binary)
	return &encoded, nil
}
