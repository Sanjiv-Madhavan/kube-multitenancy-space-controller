package shared

import (
	"context"

	"go.uber.org/zap"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func (r *Reconciler) DeleteObject(ctx context.Context, object client.Object) error {
	if err := r.Client.Delete(ctx, object); err != nil {
		return err
	}
	return nil
}

func (r *Reconciler) ProcessInProgressCondition(ctx context.Context, object ConditionsChangeObject, conditionType string, conditionStatus metav1.ConditionStatus, reason, message string) {
	r.setCondition(object, conditionType, conditionStatus, object.GetGeneration(), reason, message)
	if err := r.UpdateStatusOfResource(ctx, object); err != nil {
		r.Logger.Error("Failed to update condition", zap.Error(err))
		return
	}
}

func (r *Reconciler) ProcessReadyCondition(ctx context.Context, object ConditionsChangeObject, conditionType string, conditionStatus metav1.ConditionStatus, reason, message string) {
	r.setCondition(object, conditionType, conditionStatus, object.GetGeneration(), reason, message)
	err := r.UpdateStatusOfResource(ctx, object)
	if err != nil {
		return
	}
}

func (r *Reconciler) ProcessFailedCondition(ctx context.Context, object ConditionsChangeObject, conditionType string, conditionStatus metav1.ConditionStatus, reason, message string) {
	r.setCondition(object, conditionType, conditionStatus, object.GetGeneration(), reason, message)
	err := r.UpdateStatusOfResource(ctx, object)
	if err != nil {
		return
	}
}
