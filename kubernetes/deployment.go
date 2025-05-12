package kubernetes

import (
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func NewDeployment(name string, spec appsv1.DeploymentSpec) *appsv1.Deployment {
	return &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{Name: name},
		Spec:       spec,
	}
}

func NewNginxDeploymentSpec() appsv1.DeploymentSpec {
	nginxPodSpec := corev1.PodSpec{
		Containers: []corev1.Container{{
			Name:  "nginx",
			Image: "us.gcr.io/cf-pks-releng-a1/nginx",
			Ports: []corev1.ContainerPort{{ContainerPort: 80}},
		}},
	}
	var replicas int32
	replicas = 1
	labelMap := make(map[string]string)
	labelMap["app"] = "nginx"

	return appsv1.DeploymentSpec{
		Replicas: &replicas,
		Selector: &metav1.LabelSelector{MatchLabels: labelMap},
		Template: corev1.PodTemplateSpec{ObjectMeta: metav1.ObjectMeta{Labels: labelMap}, Spec: nginxPodSpec}}
}
