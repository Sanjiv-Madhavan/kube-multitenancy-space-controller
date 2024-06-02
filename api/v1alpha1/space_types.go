/*
Copyright 2024.

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
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// SpaceSpec defines the desired state of Space
type SpaceSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// Foo is an example field of Space. Edit space_types.go to remove/update
	// +kubebuilder:validation:Optional
	ResourceQuota corev1.ResourceQuotaSpec `json:"resourceQuota,omitempty"`

	// +kubebuilder:validation:Required
	Owners []v1.Subject `json:"owners,omitempty"`

	// +kubebuilder:validation:Optional
	AdditionalRoleBindings AdditionalRoleBindingsSpec `json:"additionalRoleBindings,omitempty"`

	// +kubebuilder:validation:Optional
	NetworkPolicies NetworkPolicy `json:"networkPolicies,omitempty"`

	// +kubebuilder:validation:Optional
	LimitRanges LimitRangesSpec `json:"limitRanges,omitempty"`

	// +kubebuilder:validation:Optional
	ServiceAccounts ServiceAccountsSpec `json:"serviceAccounts,omitempty"`

	// +kubebuilder:validation:Optional
	TemplateRef SpaceTemplateReference `json:"templateRef,omitempty"`
}

type SpaceTemplateReference struct {
	// Name of the SpaceTemplate.
	Name string `json:"name,omitempty"`
	// Kind specifies the kind of the referenced resource, which should be "SpaceTemplate".
	Kind string `json:"kind,omitempty"`
	// Group is the API group of the SpaceTemplate,  "github.sanjivmadhavan.io/v1alpha1".
	Group string `json:"group,omitempty"`
}

// SpaceStatus defines the observed state of Space
type SpaceStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// NamespaceName the name of the created underlying namespace.
	NamespaceName string `json:"namespaceName,omitempty"`
	// Conditions List of status conditions to indicate the status of Space
	Conditions []metav1.Condition `json:"conditions,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// Space is the Schema for the spaces API
type Space struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   SpaceSpec   `json:"spec,omitempty"`
	Status SpaceStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// SpaceList contains a list of Space
type SpaceList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Space `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Space{}, &SpaceList{})
}

func (s *Space) GetConditions() []metav1.Condition {
	return s.Status.Conditions
}

func (s *Space) SetConditions(conditions []metav1.Condition) {
	s.Status.Conditions = conditions
}
