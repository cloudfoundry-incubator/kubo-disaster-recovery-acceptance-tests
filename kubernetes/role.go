package kubernetes

import (
	rbac "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func NewRole(name, pspResourceName string) *rbac.Role {
	return &rbac.Role{
		ObjectMeta: metav1.ObjectMeta{Name: name},
		Rules: []rbac.PolicyRule{
			{
				APIGroups:     []string{"extensions"},
				ResourceNames: []string{pspResourceName},
				Resources:     []string{"podsecuritypolicies"},
				Verbs:         []string{"use"},
			},
		},
	}
}

func NewRoleBinding(roleName, roleBindingName, serviceAccountName string) *rbac.RoleBinding {
	return &rbac.RoleBinding{
		ObjectMeta: metav1.ObjectMeta{Name: roleBindingName},
		RoleRef: rbac.RoleRef{
			APIGroup: "rbac.authorization.k8s.io",
			Kind:     "Role",
			Name:     roleName,
		},
		Subjects: []rbac.Subject{
			{
				Kind: "ServiceAccount",
				Name: serviceAccountName,
			},
		},
	}
}
