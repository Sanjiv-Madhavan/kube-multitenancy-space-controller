package space

import (
	"context"
	"fmt"
	"strconv"

	githubsanjivmadhavaniov1alpha1 "github.com/Sanjiv-Madhavan/kube-multitenancy-space-controller/api/v1alpha1"
	"github.com/Sanjiv-Madhavan/kube-multitenancy-space-controller/internal/pkg"
	"go.uber.org/zap"
	networkingv1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

func (r *SpaceReconciler) reconcileNetworkPolicies(ctx context.Context, space *githubsanjivmadhavaniov1alpha1.Space) (err error) {

	// When creating a space in kube-multitenancy-operator, you can set the enableDefaultStrictMode parameter to true. If enabled, kube-multitenancy-operator will
	// generate a default network policy that restricts ingress communication from other spaces or namespaces. However, namespaces
	// labeled with github.sanjivmadhavan.kube-multitenancy-operator.io/role: system are exempt from this restriction.

	if space.Spec.NetworkPolicies.EnableDefaultStrictMode {
		networkPolicyName := fmt.Sprintf("github.sanjivmadhavan.kube-multitenancy-operator-custom-default-%s", space.Name)
		networkPolicySpec := newNetworkPolicyDefaultSpec()
		networkPolicy := newNetworkPolicy(networkPolicyName, space.Status.NamespaceName, networkPolicySpec)
		if err = r.syncNetworkPolicy(ctx, networkPolicy, space, networkPolicySpec); err != nil {
			r.Logger.Error("Cannot Synchronize Network policy", zap.Error(err))
			return err
		}
	}

	for idx, networkPolicyItem := range space.Spec.NetworkPolicies.Items {
		networkPolicyName := "github.sanjivmadhavan.kube-multitenancy-operator-custom-" + strconv.Itoa(idx)
		networkPolicy := newNetworkPolicy(networkPolicyName, space.Status.NamespaceName, networkPolicyItem)
		if err = r.syncNetworkPolicy(ctx, networkPolicy, space, networkPolicyItem); err != nil {
			r.Logger.Error("Cannot Synchronize Network policy", zap.Error(err))
			return err
		}
	}
	return nil
}

func newNetworkPolicy(name string, namespace string, networkPolicySpec networkingv1.NetworkPolicySpec) *networkingv1.NetworkPolicy {
	return &networkingv1.NetworkPolicy{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Spec: networkPolicySpec,
	}
}

func (r *SpaceReconciler) syncNetworkPolicy(ctx context.Context, networkPolicy *networkingv1.NetworkPolicy, space *githubsanjivmadhavaniov1alpha1.Space, spec networkingv1.NetworkPolicySpec) (err error) {
	var (
		res                            controllerutil.OperationResult
		spaceLabel, networkPolicyLabel string
	)

	if spaceLabel, err = pkg.GetTypeLabel(space); err != nil {
		return
	}

	if networkPolicyLabel, err = pkg.GetTypeLabel(networkPolicy); err != nil {
		return
	}

	res, err = controllerutil.CreateOrUpdate(ctx, r.Client, networkPolicy, func() (err error) {
		networkPolicy.SetLabels(map[string]string{
			spaceLabel:         space.Name,
			networkPolicyLabel: networkPolicy.Name,
		})
		if networkPolicy.Name != fmt.Sprintf("github.sanjivmadhavan.kube-multitenancy-operator-custom-default-%s", space.Name) {
			networkPolicy.Spec = spec
		}

		return nil
	})
	r.Logger.Info("Network Policy sync result: " + string(res) + ", name: " + networkPolicy.Name + ", namespace: " + space.Status.NamespaceName)
	r.EmitEvent(space, res, space.Name, "Ensuring network policy creation/Update", err)

	return err
}

func newNetworkPolicyDefaultSpec() networkingv1.NetworkPolicySpec {
	return networkingv1.NetworkPolicySpec{
		PodSelector: metav1.LabelSelector{MatchLabels: map[string]string{}},
		Ingress: []networkingv1.NetworkPolicyIngressRule{
			{
				From: []networkingv1.NetworkPolicyPeer{
					{
						NamespaceSelector: &metav1.LabelSelector{
							MatchLabels: map[string]string{
								"github.sanjivmadhavan.kube-multitenancy-operator.io/role": "system",
							},
						},
					},
					{
						PodSelector: &metav1.LabelSelector{
							MatchLabels: map[string]string{},
						},
					},
				},
			},
		},
	}
}

func (r *SpaceReconciler) deleteNetworkPolicies(ctx context.Context, space *githubsanjivmadhavaniov1alpha1.Space) error {
	if space.Spec.NetworkPolicies.EnableDefaultStrictMode {
		networkPolicyName := fmt.Sprintf("github.sanjivmadhavan.kube-multitenancy-operator-custom-default-%s", space.Name)
		networkPolicySpec := newNetworkPolicyDefaultSpec()
		networkPolicy := newNetworkPolicy(networkPolicyName, space.Status.NamespaceName, networkPolicySpec)
		if err := r.DeleteObject(ctx, networkPolicy); err != nil {
			r.Logger.Error("Cannot Delete Default Network policy", zap.Error(err))
			return err
		}
	}
	for idx, networkPolicyItem := range space.Spec.NetworkPolicies.Items {
		networkPolicyName := "github.sanjivmadhavan.kube-multitenancy-operator-custom-" + strconv.Itoa(idx)
		networkPolicy := newNetworkPolicy(networkPolicyName, space.Status.NamespaceName, networkPolicyItem)
		if err := r.DeleteObject(ctx, networkPolicy); err != nil {
			r.Logger.Error("Cannot Delete Network policy", zap.Error(err))
			return err
		}
	}
	return nil
}
