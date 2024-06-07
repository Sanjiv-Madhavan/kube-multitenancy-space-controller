package space

import (
	"context"
	"errors"
	"reflect"

	githubsanjivmadhavaniov1alpha1 "github.com/Sanjiv-Madhavan/kube-multitenancy-space-controller/api/v1alpha1"
	"go.uber.org/zap"
	corev1 "k8s.io/api/core/v1"
	apierrs "k8s.io/apimachinery/pkg/api/errors"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func (r *SpaceReconciler) FetchSpaceTemplate(ctx context.Context, name string) (*githubsanjivmadhavaniov1alpha1.SpaceTemplate, error) {
	spaceTemplate := &githubsanjivmadhavaniov1alpha1.SpaceTemplate{}

	if err := r.Get(ctx, client.ObjectKey{
		Name: name,
	}, spaceTemplate); err != nil {
		if apierrs.IsNotFound(err) {
			// SpaceTemplate not found, return
			r.Logger.Error("Space object not found", zap.Error(err))
			return nil, err
		}
	}
	return spaceTemplate, nil
}

func (r *SpaceReconciler) MergeResourceQuotaas(space *githubsanjivmadhavaniov1alpha1.Space, spaceTemplate *githubsanjivmadhavaniov1alpha1.SpaceTemplate) (*corev1.ResourceQuotaSpec, error) {
	resourceQuotaSpec := &corev1.ResourceQuotaSpec{}
	resourceQuotaSpec.Hard = make(corev1.ResourceList)
	switch {
	case !reflect.ValueOf(space.Spec.ResourceQuota).IsZero() && !reflect.ValueOf(spaceTemplate.Spec.ResourceQuota).IsZero():
		overrideResourceQuotas(resourceQuotaSpec, space.Spec.ResourceQuota.Hard, spaceTemplate.Spec.ResourceQuota.Hard, corev1.ResourceLimitsCPU)
		overrideResourceQuotas(resourceQuotaSpec, space.Spec.ResourceQuota.Hard, spaceTemplate.Spec.ResourceQuota.Hard, corev1.ResourceLimitsMemory)
		overrideResourceQuotas(resourceQuotaSpec, space.Spec.ResourceQuota.Hard, spaceTemplate.Spec.ResourceQuota.Hard, corev1.ResourceRequestsCPU)
		overrideResourceQuotas(resourceQuotaSpec, space.Spec.ResourceQuota.Hard, spaceTemplate.Spec.ResourceQuota.Hard, corev1.ResourceRequestsMemory)
	case reflect.ValueOf(space.Spec.ResourceQuota).IsZero() && !reflect.ValueOf(spaceTemplate.Spec.ResourceQuota).IsZero():
		space.Spec.ResourceQuota.Hard = spaceTemplate.Spec.ResourceQuota.Hard
	default:
		r.Logger.Info("merge not required both space and spacetpl resource quotas are empty")
		return nil, errors.New("merge not required both space and spacetpl resource quotas are empty")
	}
	return resourceQuotaSpec, nil
}

func overrideResourceQuotas(resourceQuotaSpec *corev1.ResourceQuotaSpec, spaceHard, spaceTemplateHard corev1.ResourceList, resourceName corev1.ResourceName) {
	if spaceTemplateHardValue, exists := spaceTemplateHard[resourceName]; exists {
		resourceQuotaSpec.Hard[resourceName] = spaceTemplateHardValue
	} else {
		resourceQuotaSpec.Hard[resourceName] = spaceHard[resourceName]
	}
}

func (r *SpaceReconciler) MergeRoleBindings(space *githubsanjivmadhavaniov1alpha1.Space, spaceTemplate *githubsanjivmadhavaniov1alpha1.SpaceTemplate) ([]githubsanjivmadhavaniov1alpha1.AdditionalRoleBindings, error) {
	mergedRoleBindings := []githubsanjivmadhavaniov1alpha1.AdditionalRoleBindings{}
	mergedRoleBindings = append(mergedRoleBindings, space.Spec.AdditionalRoleBindings...)
	for _, roleBinding := range spaceTemplate.Spec.AdditionalRoleBindings {
		if !roleBindingCollides(mergedRoleBindings, roleBinding) {
			mergedRoleBindings = append(mergedRoleBindings, roleBinding)
		}
	}
	if len(mergedRoleBindings) > 0 {
		return mergedRoleBindings, nil
	}
	return nil, errors.New("No rolebinding specified in either space or template")
}

func roleBindingCollides(mergedBindings []githubsanjivmadhavaniov1alpha1.AdditionalRoleBindings, roleBinding githubsanjivmadhavaniov1alpha1.AdditionalRoleBindings) bool {
	for _, binding := range mergedBindings {
		if binding.RoleRef.Kind == roleBinding.RoleRef.Name && binding.RoleRef.Kind == roleBinding.RoleRef.Kind {
			if reflect.DeepEqual(binding.Subjects, roleBinding.Subjects) {
				return true
			}
		}
	}
	return false
}
