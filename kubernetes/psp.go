package kubernetes

import (
	corev1 "k8s.io/api/core/v1"
	policyv1 "k8s.io/api/policy/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func NewPodSecurityPolicy(name string, spec policyv1.PodSecurityPolicySpec) *policyv1.PodSecurityPolicy {
	return &policyv1.PodSecurityPolicy{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
			Annotations: map[string]string{
				"seccomp.security.alpha.kubernetes.io/allowedProfileNames": "docker/default",
				"apparmor.security.beta.kubernetes.io/allowedProfileNames": "runtime/default",
				"seccomp.security.alpha.kubernetes.io/defaultProfileName":  "docker/default",
				"apparmor.security.beta.kubernetes.io/defaultProfileName":  "runtime/default",
			},
		},
		Spec: spec,
	}
}

func NewPodSecurityPolicySpec() policyv1.PodSecurityPolicySpec {
	var allowPrivilegedEscalation bool
	allowPrivilegedEscalation = false
	labelMap := make(map[string]string)
	labelMap["app"] = "nginx"

	return policyv1.PodSecurityPolicySpec{
		Privileged:               false,
		AllowPrivilegeEscalation: &allowPrivilegedEscalation,
		AllowedCapabilities:      []corev1.Capability{"*"},
		HostNetwork:              true,
		HostPorts: []policyv1.HostPortRange{
			{
				Min: 0,
				Max: 65535,
			},
		},
		HostIPC: true,
		HostPID: true,
		Volumes: []policyv1.FSType{
			"configMap",
			"emptyDir",
			"projected",
			"secret",
			"downwardAPI",
		},
		RunAsUser: policyv1.RunAsUserStrategyOptions{
			Rule: policyv1.RunAsUserStrategyRunAsAny,
		},
		SELinux: policyv1.SELinuxStrategyOptions{
			Rule: policyv1.SELinuxStrategyRunAsAny,
		},
		SupplementalGroups: policyv1.SupplementalGroupsStrategyOptions{
			Rule: policyv1.SupplementalGroupsStrategyRunAsAny,
		},
		FSGroup: policyv1.FSGroupStrategyOptions{
			Rule: policyv1.FSGroupStrategyRunAsAny,
		},
	}
}
