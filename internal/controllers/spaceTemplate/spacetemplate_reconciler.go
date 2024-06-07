package spaceTemplate

import (
	"context"

	githubsanjivmadhavaniov1alpha1 "github.com/Sanjiv-Madhavan/kube-multitenancy-space-controller/api/v1alpha1"
	"github.com/Sanjiv-Madhavan/kube-multitenancy-space-controller/internal/controllers/constants"
	"go.uber.org/zap"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

func (r *SpaceTemplateReconciler) reconcileDelete(ctx context.Context, spaceTemplate *githubsanjivmadhavaniov1alpha1.SpaceTemplate) (result ctrl.Result, err error) {
	if controllerutil.ContainsFinalizer(spaceTemplate, constants.SpaceFinalizer) {
		controllerutil.RemoveFinalizer(spaceTemplate, constants.SpaceFinalizer)
		if err = r.Update(ctx, spaceTemplate); err != nil {
			return ctrl.Result{}, err
		}
	}
	return ctrl.Result{}, err
}

func (r *SpaceTemplateReconciler) reconcileSpaceTemplate(ctx context.Context, spaceTemplate *githubsanjivmadhavaniov1alpha1.SpaceTemplate) (result ctrl.Result, err error) {
	r.ProcessInProgressCondition(ctx, spaceTemplate, constants.SpaceConditionCreating, metav1.ConditionUnknown, constants.SpaceCreatingReason, constants.SpaceCreatingMessage)
	r.Logger.Info("Reconciling space template.")

	if !controllerutil.ContainsFinalizer(spaceTemplate, constants.SpaceFinalizer) {
		controllerutil.AddFinalizer(spaceTemplate, constants.SpaceFinalizer)
		if err = r.Update(ctx, spaceTemplate); err != nil {
			r.Logger.Error("Unable to update finalizers in reconciliation", zap.Error(err))
			r.ProcessFailedCondition(ctx, spaceTemplate, constants.SpaceConditionFailed, metav1.ConditionFalse, constants.SpaceFailedReason, constants.SpaceSyncFailMessage)
			r.EmitEvent(spaceTemplate, controllerutil.OperationResultCreated, spaceTemplate.Name, "Ensuring SpaceTemplate creation/Update", err)
			return ctrl.Result{}, err
		}
	}

	r.ProcessReadyCondition(ctx, spaceTemplate, constants.SpaceConditionReady, metav1.ConditionTrue, constants.SpaceSyncSuccessReason, constants.SpaceSyncSuccessMessage)
	r.EmitEvent(spaceTemplate, controllerutil.OperationResultCreated, spaceTemplate.Name, "Ensuring SpaceTemplate creation/Update", err)

	return ctrl.Result{
		RequeueAfter: constants.RequeueAfter,
	}, nil
}
