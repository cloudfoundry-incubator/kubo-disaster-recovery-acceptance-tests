package testcase

import (
	"github.com/cloudfoundry-incubator/kubo-disaster-recovery-acceptance-tests/kubernetes"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/api/errors"
)

var (
	err                 error
	k8sClient           *kubernetes.Client
	DeploymentNamespace string
	nginx1Deployment    *v1.Deployment
	nginx2Deployment    *v1.Deployment
	nginx3Deployment    *v1.Deployment
)

type Deployment struct{}

func (t Deployment) Name() string {
	return "deployment"
}

func (t Deployment) BeforeBackup(config Config) {
	By("Initializing K8s client", func() {
		k8sClient, err = kubernetes.NewKubeClient()
		Expect(err).NotTo(HaveOccurred())

		nsObject, err := k8sClient.CreateNamespace("bbr")
		Expect(err).ToNot(HaveOccurred())
		DeploymentNamespace = nsObject.Name
	})

	By("Deploying workload 1 and 2", func() {
		nginx1Deployment = kubernetes.NewDeployment("nginx-1", kubernetes.NewNginxDeploymentSpec())
		nginx1Deployment, err = k8sClient.CreateDeployment(DeploymentNamespace, nginx1Deployment)
		Expect(err).ToNot(HaveOccurred())

		nginx2Deployment = kubernetes.NewDeployment("nginx-2", kubernetes.NewNginxDeploymentSpec())
		nginx2Deployment, err = k8sClient.CreateDeployment(DeploymentNamespace, nginx2Deployment)
		Expect(err).ToNot(HaveOccurred())

		err = k8sClient.WaitForDeployment(DeploymentNamespace, nginx1Deployment.Name, GinkgoWriter)
		Expect(err).NotTo(HaveOccurred())

		err = k8sClient.WaitForDeployment(DeploymentNamespace, nginx2Deployment.Name, GinkgoWriter)
		Expect(err).NotTo(HaveOccurred())
	})
}

func (t Deployment) AfterBackup(config Config) {
	By("Deleting workload 2", func() {
		err = k8sClient.DeleteDeployment(DeploymentNamespace, nginx2Deployment.Name)
		Expect(err).ToNot(HaveOccurred())
	})

	By("Deploying workload 3", func() {
		nginx3Deployment = kubernetes.NewDeployment("nginx-3", kubernetes.NewNginxDeploymentSpec())
		nginx3Deployment, err = k8sClient.CreateDeployment(DeploymentNamespace, nginx3Deployment)
		Expect(err).ToNot(HaveOccurred())

		err = k8sClient.WaitForDeployment(DeploymentNamespace, nginx3Deployment.Name, GinkgoWriter)
		Expect(err).NotTo(HaveOccurred())
	})
}

func (t Deployment) AfterRestore(config Config) {
	By("Waiting for workloads 1 and 2 to be available", func() {
		err = k8sClient.WaitForDeployment(DeploymentNamespace, nginx1Deployment.Name, GinkgoWriter)
		Expect(err).NotTo(HaveOccurred())

		err = k8sClient.WaitForDeployment(DeploymentNamespace, nginx2Deployment.Name, GinkgoWriter)
		Expect(err).NotTo(HaveOccurred())
	})

	By("Asserting that workload 3 is gone", func() {
		_, err = k8sClient.GetDeployment(DeploymentNamespace, nginx3Deployment.Name)
		Expect(err).To(HaveOccurred())

		statusErr, ok := err.(*errors.StatusError)
		Expect(ok).To(BeTrue())
		Expect(statusErr.ErrStatus.Code).To(Equal(int32(404)))
	})
}

func (t Deployment) Cleanup(config Config) {
	By("Deleting the test namespace", func() {
		err := k8sClient.DeleteNamespace(DeploymentNamespace)
		Expect(err).NotTo(HaveOccurred())
	})
}
