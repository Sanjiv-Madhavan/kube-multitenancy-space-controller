/*
Copyright 2024.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package space

import (
	"context"

	"go.uber.org/zap"
	apierrs "k8s.io/apimachinery/pkg/api/errors"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/predicate"

	githubsanjivmadhavaniov1alpha1 "github.com/Sanjiv-Madhavan/kube-multitenancy-space-controller/api/v1alpha1"
	"github.com/Sanjiv-Madhavan/kube-multitenancy-space-controller/internal/controllers/shared"
)

// SpaceReconciler reconciles a Space object
type SpaceReconciler struct {
	shared.Reconciler
}

//+kubebuilder:rbac:groups=github.sanjivmadhavan.io,resources=spaces,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=github.sanjivmadhavan.io,resources=spaces/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=github.sanjivmadhavan.io,resources=spaces/finalizers,verbs=update
//+kubebuilder:rbac:groups=core,resources=namespaces,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=core,resources=resourcequotas,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=core,resources=limitranges,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=core,resources=serviceaccounts,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=core,resources=events,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=networking.k8s.io,resources=networkpolicies,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=rbac.authorization.k8s.io,resources=rolebindings,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=rbac.authorization.k8s.io,resources=clusterroles,verbs=bind

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Space object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.17.2/pkg/reconcile
func (r *SpaceReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := r.Logger.With(zap.String("name", req.NamespacedName.String()))

	space := &githubsanjivmadhavaniov1alpha1.Space{}

	if err := r.Get(ctx, req.NamespacedName, space); err != nil {
		if !apierrs.IsNotFound(err) {
			logger.Error("Space object not found", zap.Error(err))
			return ctrl.Result{}, nil
		}
		return ctrl.Result{}, err
	}

	// Apply filters if any
	// Check if Space is about to be deleted post completion of finalizers
	if !space.DeletionTimestamp.IsZero() {
		r.reconcileDelete(ctx, space)
	}

	if space.Spec.TemplateRef.Name != "" {
		return r.reconcileSpaceFromTemplate(ctx, space)
	}

	return r.reconcileSpace(ctx, space)
}

// SetupWithManager sets up the controller with the Manager.
func (r *SpaceReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&githubsanjivmadhavaniov1alpha1.Space{}).
		WithEventFilter(predicate.GenerationChangedPredicate{}).
		WithEventFilter(ignoreDeletionPredicate()).
		Complete(r)
}

func ignoreDeletionPredicate() predicate.Predicate {
	return predicate.Funcs{
		DeleteFunc: func(de event.DeleteEvent) bool {
			return !de.DeleteStateUnknown
		},
	}
}
