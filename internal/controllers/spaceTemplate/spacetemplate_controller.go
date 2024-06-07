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

package spaceTemplate

import (
	"context"

	"go.uber.org/zap"
	apierrs "k8s.io/apimachinery/pkg/api/errors"
	ctrl "sigs.k8s.io/controller-runtime"

	githubsanjivmadhavaniov1alpha1 "github.com/Sanjiv-Madhavan/kube-multitenancy-space-controller/api/v1alpha1"
	"github.com/Sanjiv-Madhavan/kube-multitenancy-space-controller/internal/controllers/shared"
)

// SpaceTemplateReconciler reconciles a SpaceTemplate object
type SpaceTemplateReconciler struct {
	shared.Reconciler
}

//+kubebuilder:rbac:groups=github.sanjivmadhavan.io,resources=spacetemplates,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=github.sanjivmadhavan.io,resources=spacetemplates/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=github.sanjivmadhavan.io,resources=spacetemplates/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the SpaceTemplate object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.17.2/pkg/reconcile
func (r *SpaceTemplateReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {

	logger := r.Logger.With(zap.String("name", req.NamespacedName.String()))

	spaceTemplate := &githubsanjivmadhavaniov1alpha1.SpaceTemplate{}

	if err := r.Get(ctx, req.NamespacedName, spaceTemplate); err != nil {
		if !apierrs.IsNotFound(err) {
			logger.Error("Space object not found", zap.Error(err))
			return ctrl.Result{}, nil
		}
		return ctrl.Result{}, err
	}

	if !spaceTemplate.DeletionTimestamp.IsZero() {
		r.reconcileDelete(ctx, spaceTemplate)
	}
	return r.reconcileSpaceTemplate(ctx, spaceTemplate)
}

// SetupWithManager sets up the controller with the Manager.
func (r *SpaceTemplateReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&githubsanjivmadhavaniov1alpha1.SpaceTemplate{}).
		Complete(r)
}
