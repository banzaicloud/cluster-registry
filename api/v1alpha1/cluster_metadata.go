// Copyright © 2021 Banzai Cloud
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package v1alpha1

// +k8s:deepcopy-gen=true
type ClusterMetadata struct {
	Provider          string    `json:"provider,omitempty"`
	Distribution      string    `json:"distribution,omitempty"`
	KubeProxyVersions []string  `json:"kubeProxyVersions,omitempty"`
	KubeletVersions   []string  `json:"kubeletVersions,omitempty"`
	Version           string    `json:"version,omitempty"`
	Locality          *Locality `json:"locality,omitempty"`
}

// +k8s:deepcopy-gen=true
type Locality struct {
	Region  string   `json:"region,omitempty"`
	Regions []string `json:"regions,omitempty"`
	Zones   []string `json:"zones,omitempty"`
}
