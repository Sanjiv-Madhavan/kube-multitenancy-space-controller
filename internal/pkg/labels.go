package pkg

import (
	"fmt"

	githubsanjivmadhavaniov1alpha1 "github.com/Sanjiv-Madhavan/kube-multitenancy-space-controller/api/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	networkingv1 "k8s.io/api/networking/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func GetTypeLabel(t metav1.Object) (label string, err error) {
	switch v := t.(type) {
	case *githubsanjivmadhavaniov1alpha1.Space:
		return "github.sanjivmadhavan.kube-multitenancy-operator.io/space", nil
	case *corev1.Namespace:
		return "github.sanjivmadhavan.kube-multitenancy-operator.io/namespace", nil
	case *networkingv1.NetworkPolicy:
		return "github.sanjivmadhavan.kube-multitenancy-operator.io/network-policy", nil
	case *corev1.LimitRange:
		return "github.sanjivmadhavan.kube-multitenancy-operator.io/limit-range", nil
	case *corev1.ResourceQuota:
		return "github.sanjivmadhavan.kube-multitenancy-operator.io/resource-quota", nil
	case *rbacv1.RoleBinding:
		return "github.sanjivmadhavan.kube-multitenancy-operator.io/resource-quota", nil
	case *corev1.ServiceAccount:
		return "github.sanjivmadhavan.kube-multitenancy-operator.io/service-account", nil
	default:
		err = fmt.Errorf("unable to parse the resource label type: %T", v)
	}
	return "", err
}
