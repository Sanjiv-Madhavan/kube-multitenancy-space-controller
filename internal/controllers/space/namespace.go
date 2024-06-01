package space

import (
	"context"

	"go.uber.org/zap"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"

	githubsanjivmadhavaniov1alpha1 "github.com/Sanjiv-Madhavan/kube-multitenancy-space-controller/api/v1alpha1"
)

func (r *SpaceReconciler) reconcileNamespace(ctx context.Context, space *githubsanjivmadhavaniov1alpha1.Space) (err error) {
	// reconciliation of space contains mainly label syncs
	namespace := r.newNamespace(space)

	err = r.syncNamespace(ctx, space, namespace)
	if err != nil {
		r.Logger.Error("Unable to sync namespace "+namespace.Name+" to space "+space.Name, zap.Error(err))
	}

	space.Status.NamespaceName = namespace.Name
	if err = r.UpdateStatusOfResource(ctx, space); err != nil {
		r.Logger.Error("Failed to update space "+space.Name, zap.Error(err))
		return err
	}

	return nil
}

func (r *SpaceReconciler) syncNamespace(ctx context.Context, space *githubsanjivmadhavaniov1alpha1.Space, namespace *corev1.Namespace) (err error) {
	// focus on createOrUpdate call

	res, err := controllerutil.CreateOrUpdate(ctx, r.Client, namespace, func() error {
		namespace.SetLabels(map[string]string{
			"github.sanjivmadhavan.io/namespace": namespace.Name,
			"github.sanjivmadhavan.io/space":     space.Name,
		})
		return nil
	})

	r.Logger.Info("Namespace sync result: " + string(res) + "name" + namespace.Name)

	return err
}

func (r *SpaceReconciler) newNamespace(space *githubsanjivmadhavaniov1alpha1.Space) *corev1.Namespace {
	namespace := &corev1.Namespace{
		ObjectMeta: v1.ObjectMeta{
			Name:   space.Name,
			Labels: space.Labels,
		},
	}

	return namespace
}

func (r *SpaceReconciler) deleteNamespace(ctx context.Context, space *githubsanjivmadhavaniov1alpha1.Space) (err error) {
	return nil
}
