/*
Copyright 2025.

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

package v1beta1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// MapUserSpec defines the desired state of MapUser.
type MapUserSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// The User ARN to associate with the MapUser
	// +kubebuilder:validation:Pattern:=`^arn:aws:iam::[0-9]{12}:user/[a-zA-Z0-9+=,.@_/-]+$`
	UserARN string `json:"userarn"`

	// The Kubernetes groups to associate with the MapUser
	// +kubebuilder:validation:MinItems:=1
	// +kubebuilder:validation:items:MinLength:=3
	// +kubebuilder:validation:items:MaxLength:=254
	Groups []string `json:"groups"`

	// A useful description of the MapUser
	// +kubebuilder:validation:Optional
	// +kubebuilder:validation:MaxLength:=254
	Description string `json:"description"`

	// The username for the MapUser
	// +kubebuilder:validation:MaxLength:=254
	// +kubebuilder:validation:MinLength:=3
	Username string `json:"username"`
}

// MapUserStatus defines the observed state of MapUser.
type MapUserStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// MapUser is the Schema for the mapusers API.
type MapUser struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   MapUserSpec   `json:"spec,omitempty"`
	Status MapUserStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// MapUserList contains a list of MapUser.
type MapUserList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []MapUser `json:"items"`
}

func init() {
	SchemeBuilder.Register(&MapUser{}, &MapUserList{})
}
