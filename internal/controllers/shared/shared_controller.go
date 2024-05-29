package shared

import (
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
