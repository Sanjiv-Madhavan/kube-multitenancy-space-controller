package space

import (
	"context"
	"strconv"

	githubsanjivmadhavaniov1alpha1 "github.com/Sanjiv-Madhavan/kube-multitenancy-space-controller/api/v1alpha1"
	"github.com/Sanjiv-Madhavan/kube-multitenancy-space-controller/internal/pkg"
	"go.uber.org/zap"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

func (r *SpaceReconciler) reconcileLimitRanges(ctx context.Context, space *githubsanjivmadhavaniov1alpha1.Space) (err error) {
	for i, limitRangeItem := range space.Spec.LimitRanges.Items {
		limitRangeName := "github.sanjivmadhavan.kube-multitenancy-operator-custom-" + strconv.Itoa(i)
		limitRangeToConstruct := r.newLimitRange(space, limitRangeName, limitRangeItem)
		if err = r.syncLimitRange(ctx, space, limitRangeToConstruct, limitRangeItem); err != nil {
			r.Logger.Error("Unable to sync namespace "+limitRangeName+" to space "+space.Name, zap.Error(err))
			return err
		}
	}
	return nil
}

func (r *SpaceReconciler) syncLimitRange(ctx context.Context, space *githubsanjivmadhavaniov1alpha1.Space, limitRange *corev1.LimitRange, limitRangeSpec corev1.LimitRangeSpec) (err error) {

	var (
		res                         controllerutil.OperationResult
		spaceLabel, limitRangeLabel string
	)

	if spaceLabel, err = pkg.GetTypeLabel(space); err != nil {
		return
	}

	if limitRangeLabel, err = pkg.GetTypeLabel(limitRange); err != nil {
		return
	}

	res, err = controllerutil.CreateOrUpdate(ctx, r.Client, limitRange, func() error {
		limitRange.SetLabels(map[string]string{
			limitRangeLabel: limitRange.Name,
			spaceLabel:      space.Status.NamespaceName,
		})
		limitRange.Spec = limitRangeSpec
		return nil
	})

	r.Logger.Info("ResoourceQuota sync result: " + string(res) + ", name: " + limitRange.Name)
	r.EmitEvent(space, res, space.Name, "Ensuring LimitRange creation/Update", err)

	return err
}

func (r *SpaceReconciler) newLimitRange(space *githubsanjivmadhavaniov1alpha1.Space, name string, limitRangeSpec corev1.LimitRangeSpec) *corev1.LimitRange {
	return &corev1.LimitRange{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: space.Status.NamespaceName,
		},
		Spec: limitRangeSpec,
	}
}

func (r *SpaceReconciler) deleteLimitRanges(ctx context.Context, space *githubsanjivmadhavaniov1alpha1.Space) error {
	for i, limitRangeItem := range space.Spec.LimitRanges.Items {
		limitRangeName := "github.sanjivmadhavan.kube-multitenancy-operator-custom-" + strconv.Itoa(i)
		limitRangeResourceToDelete := r.newLimitRange(space, limitRangeName, limitRangeItem)
		if err := r.DeleteObject(ctx, limitRangeResourceToDelete); err != nil {
			r.Logger.Error("Failed to delete Namespace resource" + limitRangeResourceToDelete.Name)
			return err
		}
	}
	return nil
}
