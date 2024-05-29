package space

import (
	"context"

	githubsanjivmadhavaniov1alpha1 "github.com/Sanjiv-Madhavan/kube-multitenancy-space-controller/api/v1alpha1"
	"github.com/Sanjiv-Madhavan/kube-multitenancy-space-controller/internal/controllers/constants"
	"go.uber.org/zap"
	apierrs "k8s.io/apimachinery/pkg/api/errors"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

func (r *SpaceReconciler) reconcileDelete(ctx context.Context, space *githubsanjivmadhavaniov1alpha1.Space) (result ctrl.Result, err error) {
	if space.HasIgnoreUnderlyingDeletionAnnotation() {
		// Delete Finalizers alone
		if controllerutil.ContainsFinalizer(space, constants.SpaceFinalizer) {
			controllerutil.RemoveFinalizer(space, constants.SpaceFinalizer)
			if err = r.Update(ctx, space); err != nil {
				return ctrl.Result{}, nil
			}
		}
	}

	// If no annotations, remove underlying resource gracefully
	if controllerutil.ContainsFinalizer(space, constants.SpaceFinalizer) {
		controllerutil.RemoveFinalizer(space, constants.SpaceFinalizer)

		if err = r.deleteNetworkPolicies(ctx, space); client.IgnoreNotFound(err) != nil {
			return ctrl.Result{}, err
		}

		if err = r.deleteLimitRanges(ctx, space); client.IgnoreNotFound(err) != nil {
			return ctrl.Result{}, nil
		}

		if err = r.deleteAdditionalRoleBindings(ctx, space); client.IgnoreNotFound(err) != nil {
			return ctrl.Result{}, nil
		}

		if err = r.deleteOwners(ctx, space); client.IgnoreNotFound(err) != nil {
			return ctrl.Result{}, nil
		}

		if err = r.deleteResourceQuota(ctx, space); client.IgnoreNotFound(err) != nil {
			return ctrl.Result{}, err
		}

		if err = r.deleteServiceAccounts(ctx, space); client.IgnoreNotFound(err) != nil {
			return ctrl.Result{}, err
		}

		if err = r.deleteNamespace(ctx, space); client.IgnoreNotFound(err) != nil {
			return ctrl.Result{}, err
		}

		controllerutil.RemoveFinalizer(space, constants.SpaceFinalizer)

		if err = r.Update(ctx, space); err != nil {
			if !apierrs.IsNotFound(err) {
				// No error in case of resource not found
				r.Logger.Error("Unable to update resource - space", zap.Error(err))
				return ctrl.Result{}, nil
			}
			// Error in case of update
			return ctrl.Result{}, err
		}
	}

	return ctrl.Result{}, nil
}

func (r *SpaceReconciler) reconcileSpaceFromTemplate(ctx context.Context, space *githubsanjivmadhavaniov1alpha1.Space) (ctrl.Result, error) {
	return ctrl.Result{}, nil
}

func (r *SpaceReconciler) reconcileSpace(ctx context.Context, space *githubsanjivmadhavaniov1alpha1.Space) (ctrl.Result, error) {
	return ctrl.Result{}, nil
}
