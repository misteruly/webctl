/*
Copyright 2021.

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
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// WebctlSpec defines the desired state of Webctl
type WebctlSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// Foo is an example field of Webctl. Edit webctl_types.go to remove/update
	Foo  string `json:"foo,omitempty"`
	Size int32  `json:"size,omitempty"`
	//Name string `json:"name,omitempty"`
}

// WebctlStatus defines the observed state of Webctl
type WebctlStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	Nodes []string `json:"nodes,omitempty"`
	//namespace namespaceClientConfig `json:"namespace,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// Webctl is the Schema for the webctls API
type Webctl struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   WebctlSpec   `json:"spec,omitempty"`
	Status WebctlStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// WebctlList contains a list of Webctl
type WebctlList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Webctl `json:"items"`
}

//type namespaceClientConfig struct {
//	namespace string
//}
//
//func (c namespaceClientConfig) Namespace() (string, bool, error) {
//	return c.namespace, false, nil
//}

func init() {
	SchemeBuilder.Register(&Webctl{}, &WebctlList{})
}
