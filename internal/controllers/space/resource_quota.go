package space

import (
	"context"

	githubsanjivmadhavaniov1alpha1 "github.com/Sanjiv-Madhavan/kube-multitenancy-space-controller/api/v1alpha1"
	"github.com/Sanjiv-Madhavan/kube-multitenancy-space-controller/internal/pkg"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

func (r *SpaceReconciler) reconcileResourceQuota(ctx context.Context, space *githubsanjivmadhavaniov1alpha1.Space) (err error) {
	// check finalizer and attach one if needed - not needed because it does not own any resource
	resourceQuota := r.newResourceQuota(space)
	r.syncResourceQuota(ctx, space, resourceQuota)

	return err
}

func (r *SpaceReconciler) syncResourceQuota(ctx context.Context, space *githubsanjivmadhavaniov1alpha1.Space, resourceQuota *corev1.ResourceQuota) (err error) {
	var (
		resourceQuotaLabel, spaceLabel string
	)

	if resourceQuotaLabel, err = pkg.GetTypeLabel(resourceQuota); err != nil {
		return
	}

	if spaceLabel, err = pkg.GetTypeLabel(space); err != nil {
		return
	}

	res, err := controllerutil.CreateOrUpdate(ctx, r.Client, resourceQuota, func() error {
		resourceQuota.SetLabels(map[string]string{
			resourceQuotaLabel: resourceQuota.Name,
			spaceLabel:         space.Status.NamespaceName,
		})
		resourceQuota.Spec = space.Spec.ResourceQuota
		return nil
	})

	r.Logger.Info("ResoourceQuota sync result: " + string(res) + ", name: " + resourceQuota.Name)
	return err
}

func (r *SpaceReconciler) newResourceQuota(space *githubsanjivmadhavaniov1alpha1.Space) *corev1.ResourceQuota {
	return &corev1.ResourceQuota{
		ObjectMeta: metav1.ObjectMeta{
			Name:   space.Name,
			Labels: space.Labels,
		},
		Spec: space.Spec.ResourceQuota,
	}
}

func (r *SpaceReconciler) deleteResourceQuota(ctx context.Context, space *githubsanjivmadhavaniov1alpha1.Space) (err error) {
	resourceQuotaToDelete := r.newResourceQuota(space)
	if err = r.DeleteObject(ctx, resourceQuotaToDelete); err != nil {
		r.Logger.Error("Failed to delete Namespace resource" + resourceQuotaToDelete.Name)
		return err
	}
	return nil
}
