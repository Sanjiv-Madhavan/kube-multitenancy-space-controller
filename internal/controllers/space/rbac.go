package space

import (
	"context"

	githubsanjivmadhavaniov1alpha1 "github.com/Sanjiv-Madhavan/kube-multitenancy-space-controller/api/v1alpha1"
)

func (r *SpaceReconciler) deleteOwners(ctx context.Context, space *githubsanjivmadhavaniov1alpha1.Space) (err error) {
	return nil
}

func (r *SpaceReconciler) deleteAdditionalRoleBindings(ctx context.Context, space *githubsanjivmadhavaniov1alpha1.Space) (err error) {
	return nil
}