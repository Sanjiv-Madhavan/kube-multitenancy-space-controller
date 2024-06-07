package space

import (
	"context"

	"go.uber.org/zap"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"

	githubsanjivmadhavaniov1alpha1 "github.com/Sanjiv-Madhavan/kube-multitenancy-space-controller/api/v1alpha1"
	"github.com/Sanjiv-Madhavan/kube-multitenancy-space-controller/internal/pkg"
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

	var namespaceLabel, spaceLabel string

	if namespaceLabel, err = pkg.GetTypeLabel(namespace); err != nil {
		return
	}

	if spaceLabel, err = pkg.GetTypeLabel(space); err != nil {
		return
	}

	res, err := controllerutil.CreateOrUpdate(ctx, r.Client, namespace, func() error {
		namespace.SetLabels(map[string]string{
			namespaceLabel: namespace.Name,
			spaceLabel:     space.Name,
		})
		return nil
	})

	r.Logger.Info("Namespace sync result: " + string(res) + "name" + namespace.Name)
	r.EmitEvent(space, res, space.Name, "Ensuring Namespace creation/Update", err)

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
	namespaceResourceToDelete := r.newNamespace(space)

	if err = r.DeleteObject(ctx, namespaceResourceToDelete); err != nil {
		r.Logger.Error("Failed to delete Namespace resource" + namespaceResourceToDelete.Name)
		return err
	}
	return nil
}
