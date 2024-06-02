package shared

import (
	"context"

	"go.uber.org/zap"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/record"
	"k8s.io/utils/clock"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type Reconciler struct {
	client.Client
	Scheme        *runtime.Scheme
	EventRecorder record.EventRecorder
	Logger        *zap.Logger
}

func (r *Reconciler) UpdateStatusOfResource(ctx context.Context, obj client.Object) error {
	if err := r.Client.Status().Update(ctx, obj); err != nil {
		r.Logger.Error("Failed to update status of the object: "+obj.GetName(), zap.Error(err))
		return err
	}
	return nil
}

// Clock is defined as a package var so it can be stubbed out during tests.
var Clock clock.Clock = clock.RealClock{}

type ConditionsChangeObject interface {
	client.Object
	metav1.Object
	GetConditions() []metav1.Condition
	SetConditions([]metav1.Condition)
}

func (r *Reconciler) setCondition(object ConditionsChangeObject, conditionType string, conditionStatus metav1.ConditionStatus, ObservedGeneration int64, reason string, message string) {
	newCondition := &metav1.Condition{
		Type:               conditionType,
		Status:             conditionStatus,
		ObservedGeneration: ObservedGeneration,
		Reason:             reason,
		Message:            message,
	}

	currentTime := metav1.NewTime(Clock.Now())
	newCondition.LastTransitionTime = currentTime

	// Update the condition
	existingConditions := object.GetConditions()
	for index, condition := range existingConditions {
		if condition.Type != conditionType {
			continue
		}

		if condition.Status == conditionStatus {
			newCondition.LastTransitionTime = condition.LastTransitionTime
		} else {
			changeObjectLogger := r.Logger.Named(object.GetName())
			changeObjectLogger.Info("Found status change for Space condition, setting lastTransitionTime to: " + currentTime.String())
		}

		existingConditions[index] = *newCondition

		return
	}
}
