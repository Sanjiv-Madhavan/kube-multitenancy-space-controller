package shared

import (
	"context"

	"go.uber.org/zap"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/record"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type Reconciler struct {
	client.Client
	Scheme        runtime.Scheme
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
