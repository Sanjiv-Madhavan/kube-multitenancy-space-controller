package v1alpha1

import (
	v1 "k8s.io/api/networking/v1"
)

type NetworkPolicy struct {
	EnableDefaultStrictMode bool                   `json:"enableDefaultStrictMode,omitempty"`
	Items                   []v1.NetworkPolicySpec `json:"items,omitempty"`
}
