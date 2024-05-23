package v1alpha1

import (
	v1 "k8s.io/api/rbac/v1"
)

type AdditionalRoleBindings struct {
	Subjects []v1.Subject `json:"subjects,omitempty"`
	RoleRef  v1.RoleRef   `json:"roleRef,omitempty"`
}

type AdditionalRoleBindingsSpec []AdditionalRoleBindings
