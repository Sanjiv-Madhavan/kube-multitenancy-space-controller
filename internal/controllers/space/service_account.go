package space

import (
	"context"

	githubsanjivmadhavaniov1alpha1 "github.com/Sanjiv-Madhavan/kube-multitenancy-space-controller/api/v1alpha1"
	"github.com/Sanjiv-Madhavan/kube-multitenancy-space-controller/internal/pkg"
	"go.uber.org/zap"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

func (r *SpaceReconciler) reconcileServiceAccounts(ctx context.Context, space *githubsanjivmadhavaniov1alpha1.Space) (err error) {
	for _, serviceAccount := range space.Spec.ServiceAccounts.Items {
		serviceAccount := newServiceAccount(serviceAccount.Name, space.Status.NamespaceName, serviceAccount.Annotations)
		err = r.syncServiceAccount(ctx, serviceAccount, space, serviceAccount.Annotations)

		if err != nil {
			r.Logger.Error("Cannot Synchronize Service Account", zap.Error(err))
			return err
		}
	}

	return nil
}

func (r *SpaceReconciler) syncServiceAccount(ctx context.Context, serviceAccount *corev1.ServiceAccount, space *githubsanjivmadhavaniov1alpha1.Space, annotations githubsanjivmadhavaniov1alpha1.Annotations) (err error) {
	var (
		res                             controllerutil.OperationResult
		spaceLabel, serviceAccountLabel string
	)

	if spaceLabel, err = pkg.GetTypeLabel(space); err != nil {
		return
	}

	if serviceAccountLabel, err = pkg.GetTypeLabel(serviceAccount); err != nil {
		return
	}

	res, err = controllerutil.CreateOrUpdate(ctx, r.Client, serviceAccount, func() (err error) {
		serviceAccount.SetLabels(map[string]string{
			spaceLabel:          space.Name,
			serviceAccountLabel: serviceAccount.Name,
		})
		serviceAccount.SetAnnotations(annotations)

		return nil
	})
	r.Logger.Info("ServiceAccount sync result: " + string(res) + ", name: " + serviceAccount.Name + ", namespace: " + space.Status.NamespaceName)

	return err
}

func newServiceAccount(name, namespace string, annotations githubsanjivmadhavaniov1alpha1.Annotations) *corev1.ServiceAccount {
	return &corev1.ServiceAccount{
		ObjectMeta: metav1.ObjectMeta{
			Name:        name,
			Namespace:   namespace,
			Annotations: annotations,
		},
	}
}

func (r *SpaceReconciler) deleteServiceAccounts(ctx context.Context, space *githubsanjivmadhavaniov1alpha1.Space) (err error) {
	for _, serviceAccount := range space.Spec.ServiceAccounts.Items {
		serviceAccount := newServiceAccount(serviceAccount.Name, space.Status.NamespaceName, serviceAccount.Annotations)
		err = r.DeleteObject(ctx, serviceAccount)

		if err != nil {
			r.Logger.Error("Cannot Synchronize Service Account", zap.Error(err))
			return err
		}
	}
	return nil
}
