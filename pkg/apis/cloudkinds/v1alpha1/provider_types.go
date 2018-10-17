/*
Copyright 2018 The Kubernetes Authors.

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

// ProviderSpec defines the desired state of Provider
type ProviderSpec struct {
	// WebHook is the URL to call when a resource should be acted upon by the provider.
	WebHook string `json:"webhook"`

	// SupportedKinds for the provider.
	// Example:
	// - MySQLDatabase
	// - PostgreSQLServer
	Kinds []string `json:"kinds"`
}

// ProviderStatus defines the observed state of Provider
type ProviderStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// Provider is the Schema for the providers API
// +k8s:openapi-gen=true
type Provider struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ProviderSpec   `json:"spec,omitempty"`
	Status ProviderStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ProviderList contains a list of Provider
type ProviderList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Provider `json:"items"`
}

type ResourceAction string

const (
	ResourceCreated ResourceAction = "create"
	ResourceUpdated ResourceAction = "update"
	ResourceDeleted ResourceAction = "delete"
)

type ResourceReference struct {
	metav1.TypeMeta `json:",inline"`
	Name            string `json:"name"`
	Namespace       string `json:"namespace"`
}

type ResourceEvent struct {
	Action   ResourceAction    `json:"action"`
	Resource ResourceReference `json:"resource"`
}

func init() {
	SchemeBuilder.Register(&Provider{}, &ProviderList{})
}
