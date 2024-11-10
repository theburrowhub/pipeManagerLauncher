package namespace

import (
	"context"
	"fmt"

	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// roleBinding creates role bindings for the given roles and binds them to the given Service Account
func roleBinding(client *kubernetes.Clientset, namespace string, saName string, roleNames []string) error {
	for _, roleName := range roleNames {
		bindingName := fmt.Sprintf("%.59s-%.3s-binding", roleName, saName)

		roleBinding := &rbacv1.RoleBinding{
			ObjectMeta: metav1.ObjectMeta{
				Name:      bindingName,
				Namespace: namespace,
			},
			Subjects: []rbacv1.Subject{
				{
					Kind:      "ServiceAccount",
					Name:      saName,
					Namespace: namespace,
				},
			},
			RoleRef: rbacv1.RoleRef{
				APIGroup: "rbac.authorization.k8s.io",
				Kind:     "Role",
				Name:     roleName,
			},
		}

		_, err := client.RbacV1().RoleBindings(namespace).Get(context.TODO(), roleBinding.Name, metav1.GetOptions{})
		if err != nil {
			_, err = client.RbacV1().RoleBindings(namespace).Create(context.TODO(), roleBinding, metav1.CreateOptions{})
			if err != nil {
				return fmt.Errorf("failed to create role binding for role %s: %w", roleName, err)
			}
		}
	}

	return nil
}
