package testcases

import (
	"k8s.io/client-go/kubernetes"
	tappsv1 "k8s.io/client-go/kubernetes/typed/apps/v1"
	restclient "k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	appsv1 "k8s.io/api/apps/v1"
	. "github.com/onsi/gomega"
	. "github.com/onsi/ginkgo"
	"github.com/satori/go.uuid"
	"k8s.io/api/apps/v1"
	"time"
	"k8s.io/apimachinery/pkg/watch"
	"fmt"
	"net/http"
	"k8s.io/apimachinery/pkg/api/errors"
	. "github.com/cloudfoundry-incubator/kubo-disaster-recovery-acceptance-tests/helpers"
)

type KuboTestCase struct{
}

var (
	err            error
	k8s            kubernetes.Interface
	namespace      string
	deploymentApi  tappsv1.DeploymentInterface
	nginx1Deployment *v1.Deployment
	nginx2Deployment *v1.Deployment
	nginx3Deployment *v1.Deployment

)

func (t KuboTestCase) Name() string {
	return "kubo_test_case"
}

func (t KuboTestCase) BeforeBackup(config Config) {
	By("Initializing K8s client", func() {
		k8s, err = newKubeClient()
		Expect(err).NotTo(HaveOccurred())

		nsObject, err := createTestNamespace(k8s, "bbr")
		Expect(err).ToNot(HaveOccurred())
		namespace = nsObject.Name
		deploymentApi = k8s.AppsV1().Deployments(namespace)
	})

	By("Deploying workload 1 and 2", func() {
		nginx1Deployment = newDeployment("nginx-1", getNginxDeploymentSpec())
		nginx1Deployment, err = deploymentApi.Create(nginx1Deployment)
		Expect(err).ToNot(HaveOccurred())

		nginx2Deployment = newDeployment("nginx-2", getNginxDeploymentSpec())
		nginx2Deployment, err = deploymentApi.Create(nginx2Deployment)
		Expect(err).ToNot(HaveOccurred())

		err = waitForDeployment(deploymentApi, namespace, nginx1Deployment.Name)
		Expect(err).NotTo(HaveOccurred())
		err = waitForDeployment(deploymentApi, namespace, nginx2Deployment.Name)
		Expect(err).NotTo(HaveOccurred())
	})
}

func (t KuboTestCase) AfterBackup(config Config) {
	err = deploymentApi.Delete(nginx2Deployment.Name, &metav1.DeleteOptions{})
	Expect(err).ToNot(HaveOccurred())

	nginx3Deployment = newDeployment("nginx-3", getNginxDeploymentSpec())
	nginx3Deployment, err = k8s.AppsV1().Deployments(namespace).Create(nginx3Deployment)
	Expect(err).ToNot(HaveOccurred())
	err = waitForDeployment(deploymentApi, namespace, nginx3Deployment.Name)
	Expect(err).NotTo(HaveOccurred())
}

func (t KuboTestCase) AfterRestore(config Config) {
	By("Waiting for API to be available", func() {
		Eventually(func() bool {
			var status int
			k8s.CoreV1().RESTClient().Get().RequestURI("/healthz").Do().StatusCode(&status)
			if status == http.StatusOK {
				return true
			}
			return false
		}, "60s", "5s").Should(BeTrue())

	})

	By("Waiting for workloads 1 and 2 to be available", func() {
		err = waitForDeployment(deploymentApi, namespace, nginx1Deployment.Name)
		Expect(err).NotTo(HaveOccurred())
		err = waitForDeployment(deploymentApi, namespace, nginx2Deployment.Name)
		Expect(err).NotTo(HaveOccurred())
	})

	By("Asserting that workload 3 is gone", func() {
		_, err = deploymentApi.Get(nginx3Deployment.Name, metav1.GetOptions{})
		Expect(err).To(HaveOccurred())
		statusErr, ok := err.(*errors.StatusError)
		Expect(ok).To(BeTrue())
		Expect(statusErr.ErrStatus.Code).To(Equal(int32(404)))
	})

	By("Waiting for system workloads", func() {
		expectedSelector := []string{"kube-dns", "heapster", "kubernetes-dashboard", "influxdb"}
		runner := NewKubectlRunner()

		systemDeploymentApi := k8s.AppsV1().Deployments("kube-system")
		for _, selector := range expectedSelector {
			deployment := runner.GetResourceNameBySelector("kube-system", "deployment", "k8s-app="+selector)
			err = waitForDeployment(systemDeploymentApi, "kube-system", deployment)
			Expect(err).NotTo(HaveOccurred())
		}
	})
}

func (t KuboTestCase) Cleanup(config Config) {}

func newKubeClient() (kubernetes.Interface, error) {
	config, err := readKubeConfig()
	if err != nil {
		return nil, err
	}
	return kubernetes.NewForConfig(config)
}

func readKubeConfig() (*restclient.Config, error) {
	return clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
		clientcmd.NewDefaultClientConfigLoadingRules(),
		&clientcmd.ConfigOverrides{},
	).ClientConfig()
}

func createTestNamespace(k8s kubernetes.Interface, prefix string) (*corev1.Namespace, error) {
	name := prefix + "-" + uuid.NewV4().String()
	labels := make(map[string]string)
	labels["test"] = prefix
	namespaceObject := corev1.Namespace{ObjectMeta: metav1.ObjectMeta{Name: name, Labels: labels}}
	return k8s.CoreV1().Namespaces().Create(&namespaceObject)
}

func newDeployment(name string, spec appsv1.DeploymentSpec) *appsv1.Deployment {
	return &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{Name: name},
		Spec:       spec,
	}
}

func getNginxDeploymentSpec() appsv1.DeploymentSpec {
	nginxPodSpec := corev1.PodSpec{
		Containers: []corev1.Container{{
			Name:  "nginx",
			Image: "nginx",
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

func waitForDeployment(deploymentAPI tappsv1.DeploymentInterface, namespace string, deploymentName string) error {
	w, err := deploymentAPI.Watch(metav1.ListOptions{})
	if err != nil {
		return err
	}

	_, err = watch.Until(1*time.Minute, w, func(event watch.Event) (bool, error) {
		deployment, ok := event.Object.(*appsv1.Deployment)
		if !ok {
			return false, fmt.Errorf("Expected `%#v` to be of type appsv1.Deployment", event.Object)
		}

		if deployment.Name == deploymentName {
			if deployment.Status.AvailableReplicas == deployment.Status.UpdatedReplicas {
				return true, nil
			}
			fmt.Fprintf(GinkgoWriter, "Expected %d to be equal to %d\n", deployment.Status.AvailableReplicas, deployment.Status.UpdatedReplicas)
		}

		return false, nil
	})

	if err != nil {
		return fmt.Errorf("Deployment `%s` did not finish rolling out with error: %s", deploymentName, err)
	}

	return nil
}