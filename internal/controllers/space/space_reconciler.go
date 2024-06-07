package space

import (
	"context"
	"reflect"

	githubsanjivmadhavaniov1alpha1 "github.com/Sanjiv-Madhavan/kube-multitenancy-space-controller/api/v1alpha1"
	"github.com/Sanjiv-Madhavan/kube-multitenancy-space-controller/internal/controllers/constants"
	"go.uber.org/zap"
	apierrs "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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

func (r *SpaceReconciler) reconcileSpace(ctx context.Context, space *githubsanjivmadhavaniov1alpha1.Space) (ctrl.Result, error) {

	// Check finalizers
	if !controllerutil.ContainsFinalizer(space, constants.SpaceFinalizer) {
		controllerutil.AddFinalizer(space, constants.SpaceFinalizer)
		// Don't forget to update the resource
		if err := r.Update(ctx, space); err != nil {
			r.Logger.Error("Unable to update finalizers in reconciliation", zap.Error(err))
			r.EmitEvent(space, controllerutil.OperationResultNone, space.Name, "Space creation/update Failure", err)
			return ctrl.Result{}, err
		}
	}

	// Space implements conditionsChageObject interface coz of typeMeta and objectMeta in struct
	r.ProcessInProgressCondition(ctx, space, constants.SpaceConditionCreating, metav1.ConditionUnknown, constants.SpaceCreatingReason, constants.SpaceCreatingMessage)
	// implement metrics in future

	r.Logger.Info("Reconciling Namespace for space.")

	if err := r.reconcileNamespace(ctx, space); err != nil {
		r.Logger.Error("Failed to reconcile the namespace for the space", zap.Error(err))
		r.ProcessFailedCondition(ctx, space, constants.SpaceConditionFailed, metav1.ConditionFalse, constants.SpaceFailedReason, constants.SpaceSyncFailMessage)
		// set metrics

		return ctrl.Result{}, err
	}

	// ReconcileResourcQuota - check if resource quota is set
	resourceQuotaSpecValue := reflect.ValueOf(space.Spec.ResourceQuota)
	if !resourceQuotaSpecValue.IsZero() {
		r.Logger.Info("Reconciling resource quota")
		if err := r.reconcileResourceQuota(ctx, space); err != nil {
			r.Logger.Error("Failed reconciling resource quota", zap.Error(err))
			r.ProcessFailedCondition(ctx, space, constants.SpaceConditionFailed, metav1.ConditionFalse, constants.SpaceFailedReason, constants.SpaceSyncFailMessage)
			return ctrl.Result{}, err
		}
	}

	limitRangeSpecValue := reflect.ValueOf(space.Spec.LimitRanges)
	if !limitRangeSpecValue.IsZero() {
		r.Logger.Info("Reconciling limit ranges")
		if err := r.reconcileLimitRanges(ctx, space); err != nil {
			r.Logger.Error("Failed reconciling Limit Ranges", zap.Error(err))
			r.ProcessFailedCondition(ctx, space, constants.SpaceConditionFailed, metav1.ConditionFalse, constants.SpaceFailedReason, constants.SpaceSyncFailMessage)
			return ctrl.Result{}, err
		}
	}

	ownerRoleBindingSpecValue := reflect.ValueOf(space.Spec.Owners)
	if !ownerRoleBindingSpecValue.IsZero() {
		r.Logger.Info("Reconciling Owner Role Binding for space")
		if err := r.reconcileOwners(ctx, space); err != nil {
			r.Logger.Error("Failed reconciling Owner Role Bindings", zap.Error(err))
			r.ProcessFailedCondition(ctx, space, constants.SpaceConditionFailed, metav1.ConditionFalse, constants.SpaceFailedReason, constants.SpaceSyncFailMessage)
			return ctrl.Result{}, err
		}
	}

	additionalRoleBindings := reflect.ValueOf(space.Spec.AdditionalRoleBindings)
	if !additionalRoleBindings.IsZero() {
		r.Logger.Info("Reconciling Additional RoleBindings")
		if err := r.reconcileAdditionalRoleBindings(ctx, space); err != nil {
			r.Logger.Error("Failed reconciling Additional Role Bindings", zap.Error(err))
			r.ProcessFailedCondition(ctx, space, constants.SpaceConditionFailed, metav1.ConditionFalse, constants.SpaceFailedReason, constants.SpaceSyncFailMessage)
			return ctrl.Result{}, err
		}
	}

	networkPoliciesSpecValue := reflect.ValueOf(space.Spec.NetworkPolicies)
	if !networkPoliciesSpecValue.IsZero() {
		r.Logger.Info("Reconciling Network Policies")
		if err := r.reconcileNetworkPolicies(ctx, space); err != nil {
			r.Logger.Error("Failed reconciling Network Policies", zap.Error(err))
			r.ProcessFailedCondition(ctx, space, constants.SpaceConditionFailed, metav1.ConditionFalse, constants.SpaceFailedReason, constants.SpaceSyncFailMessage)
			return ctrl.Result{}, err
		}
	}

	serviceAccountSpecValue := reflect.ValueOf(space.Spec.ServiceAccounts)
	if !serviceAccountSpecValue.IsZero() {
		r.Logger.Info("Reconciling Service Accounts")
		if err := r.reconcileServiceAccounts(ctx, space); err != nil {
			r.Logger.Error("Failed reconciling Service Accounts", zap.Error(err))
			r.ProcessFailedCondition(ctx, space, constants.SpaceConditionFailed, metav1.ConditionFalse, constants.SpaceFailedReason, constants.SpaceSyncFailMessage)
			return ctrl.Result{}, err
		}
	}

	r.ProcessReadyCondition(ctx, space, constants.SpaceConditionReady, metav1.ConditionTrue, constants.SpaceSyncSuccessReason, constants.SpaceSyncSuccessMessage)
	r.EmitEvent(space, controllerutil.OperationResultCreated, space.Name, "Ensuring space creation/update Failure", nil)

	return ctrl.Result{
		RequeueAfter: constants.RequeueAfter,
	}, nil
}

func (r *SpaceReconciler) reconcileSpaceFromTemplate(ctx context.Context, space *githubsanjivmadhavaniov1alpha1.Space) (ctrl.Result, error) {

	spaceTemplate, err := r.FetchSpaceTemplate(ctx, space.Spec.TemplateRef.Name)
	if err != nil {
		r.Logger.Info("SpaceTemplate: " + space.Spec.TemplateRef.Name + "SpaceTemplate does not exist")
		return ctrl.Result{}, err
	}

	if res, err := r.MergeResourceQuotaas(space, spaceTemplate); err != nil {
		space.Spec.ResourceQuota = *res
	}

	if res, err := r.MergeRoleBindings(space, spaceTemplate); err != nil {
		space.Spec.AdditionalRoleBindings = res
	}

	if reflect.ValueOf(space.Spec.NetworkPolicies).IsZero() {
		space.Spec.NetworkPolicies = spaceTemplate.Spec.NetworkPolicies
	}

	if reflect.ValueOf(space.Spec.LimitRanges).IsZero() {
		space.Spec.LimitRanges = spaceTemplate.Spec.LimitRanges
	}

	r.Logger.Info("Reconciling space: " + string(space.Name) + "from template: " + spaceTemplate.Name)

	return r.reconcileSpace(ctx, space)
}
