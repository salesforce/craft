// Copyright (c) 2020, salesforce.com, inc.
// All rights reserved.
// SPDX-License-Identifier: BSD-3-Clause
// For full license text, see LICENSE.txt file in the repo root or https://opensource.org/licenses/BSD-3-Clause

/*

Licensed under the Apache License, .Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package {{ .Version }}

import (
	meta{{ .Version }} "k8s.io/apimachinery/pkg/apis/meta/v1"
	corev1 "k8s.io/api/core/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// {{ .Resource }}Spec defines the desired state of {{ .Resource }}

type Pod struct {
	Name string `json:"name,omitempty"`
	Namespace string `json:"namespace,omitempty"`
	Type string `json:"type,omitempty"`
}

// {{ .Resource }}Status defines the observed state of {{ .Resource }}
type {{ .Resource }}Status struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	StatusPayload string `json:"statusPayload,omitempty"`
	Pod Pod `json:"pod,omitempty"`
	State  string `json:"state,omitempty"`
	Message string `json:"message,omitempty"`
	Terminated *corev1.ContainerStateTerminated `json:"terminated,omitempty"`
}

// +kubebuilder:object:root=true

// {{ .Resource }} is the Schema for the {{ .Resource }}s API
type {{ .Resource }} struct {
	meta{{ .Version }}.TypeMeta   `json:",inline"`
	meta{{ .Version }}.ObjectMeta `json:"metadata,omitempty"`

	Spec   Root   `json:"spec,omitempty"`
	Status {{ .Resource }}Status `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// {{ .Resource }}List contains a list of {{ .Resource }}
type {{ .Resource }}List struct {
	meta{{ .Version }}.TypeMeta `json:",inline"`
	meta{{ .Version }}.ListMeta `json:"metadata,omitempty"`
	Items           []{{ .Resource }} `json:"items"`
}

func init() {
	SchemeBuilder.Register(&{{ .Resource }}{}, &{{ .Resource }}List{})
}
