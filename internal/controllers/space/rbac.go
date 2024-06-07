package space

import (
	"context"

	githubsanjivmadhavaniov1alpha1 "github.com/Sanjiv-Madhavan/kube-multitenancy-space-controller/api/v1alpha1"
	"github.com/Sanjiv-Madhavan/kube-multitenancy-space-controller/internal/pkg"
	"go.uber.org/zap"
	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

func (r *SpaceReconciler) reconcileOwners(ctx context.Context, space *githubsanjivmadhavaniov1alpha1.Space) (err error) {
	roleBindingName := space.Name + "-owner"
	roleRef := rbacv1.RoleRef{
		Kind:     "ClusterRole",
		Name:     "admin", // Default admin role
		APIGroup: rbacv1.GroupName,
	}
	roleBinding := r.newRoleBinding(space, roleBindingName, roleRef, space.Spec.Owners)
	err = r.syncRoleBinding(ctx, space, roleBindingName, roleBinding, roleRef, roleBinding.Subjects)
	return err
}

func (r *SpaceReconciler) reconcileAdditionalRoleBindings(ctx context.Context, space *githubsanjivmadhavaniov1alpha1.Space) (err error) {
	for _, addOnRoleBinding := range space.Spec.AdditionalRoleBindings {
		roleBindingName := space.Name + "-" + addOnRoleBinding.RoleRef.Name
		additionalRoleBindingToConstruct := r.newRoleBinding(space, roleBindingName, addOnRoleBinding.RoleRef, addOnRoleBinding.Subjects)

		err = r.syncRoleBinding(ctx, space, roleBindingName, additionalRoleBindingToConstruct, addOnRoleBinding.RoleRef, addOnRoleBinding.Subjects)
	}

	return err
}

func (r *SpaceReconciler) syncRoleBinding(ctx context.Context, space *githubsanjivmadhavaniov1alpha1.Space, rolebindingName string, roleBinding *rbacv1.RoleBinding, roleRef rbacv1.RoleRef, subjects []rbacv1.Subject) (err error) {
	var (
		res                          controllerutil.OperationResult
		spaceLabel, roleBindingLabel string
	)

	if spaceLabel, err = pkg.GetTypeLabel(space); err != nil {
		return
	}

	if roleBindingLabel, err = pkg.GetTypeLabel(roleBinding); err != nil {
		return
	}

	res, err = controllerutil.CreateOrUpdate(ctx, r.Client, roleBinding, func() error {
		roleBinding.SetLabels(map[string]string{
			spaceLabel:       space.Name,
			roleBindingLabel: rolebindingName,
		})
		roleBinding.RoleRef = roleRef
		roleBinding.Subjects = subjects

		return nil
	})
	r.Logger.Info("Rolebinding sync result: " + string(res) + ", name: " + rolebindingName)
	r.EmitEvent(space, res, space.Name, "Ensuring RoleBinding creation/Update", err)

	return err
}

func (r *SpaceReconciler) newRoleBinding(space *githubsanjivmadhavaniov1alpha1.Space, name string, roleRef rbacv1.RoleRef, subjects []rbacv1.Subject) *rbacv1.RoleBinding {
	return &rbacv1.RoleBinding{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: space.Status.NamespaceName,
		},
		RoleRef:  roleRef,
		Subjects: subjects,
	}
}

func (r *SpaceReconciler) deleteOwners(ctx context.Context, space *githubsanjivmadhavaniov1alpha1.Space) (err error) {
	roleBindingName := space.Name + "-owner"
	roleRef := rbacv1.RoleRef{}

	roleBinding := r.newRoleBinding(space, roleBindingName, roleRef, space.Spec.Owners)
	err = r.DeleteObject(ctx, roleBinding)

	return err
}

func (r *SpaceReconciler) deleteAdditionalRoleBindings(ctx context.Context, space *githubsanjivmadhavaniov1alpha1.Space) (err error) {
	for _, addOnRoleBinding := range space.Spec.AdditionalRoleBindings {
		roleBindingName := space.Name + "-" + addOnRoleBinding.RoleRef.Name
		additionalRoleBindingToConstruct := r.newRoleBinding(space, roleBindingName, addOnRoleBinding.RoleRef, addOnRoleBinding.Subjects)

		if err = r.DeleteObject(ctx, additionalRoleBindingToConstruct); err != nil {
			r.Logger.Error("Cannot Delete Addon Role Binding", zap.Error(err))
			return err
		}
	}
	return nil
}
