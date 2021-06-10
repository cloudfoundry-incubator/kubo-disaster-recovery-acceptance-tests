package testcase

import (
	"time"

	"github.com/cloudfoundry-incubator/kubo-disaster-recovery-acceptance-tests/kubernetes"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/api/errors"
)

type Deployment struct {
	k8sClient             *kubernetes.Client
	namespace             string
	nginx1Deployment      *v1.Deployment
	nginx2Deployment      *v1.Deployment
	nginx3Deployment      *v1.Deployment
	timeout               time.Duration
	serviceAccountName    string
	podSecurityPolicyName string
	roleName              string
	roleBindingName       string
}

func NewDeployment() *Deployment {
	return &Deployment{}
}

func (t *Deployment) Name() string {
	return "deployment"
}

func (t *Deployment) BeforeBackup(config Config) {
	By("Initializing K8s client", func() {
		var err error
		t.k8sClient, err = kubernetes.NewKubeClient()
		Expect(err).NotTo(HaveOccurred())

		nsObject, err := t.k8sClient.CreateNamespace("bbr")
		Expect(err).ToNot(HaveOccurred())
		t.namespace = nsObject.Name

		t.timeout = 20 * time.Minute
	})

	By("Creating a service account with PSP privileges", func() {
		t.serviceAccountName = "deployment"
		_, err := t.k8sClient.CreateServiceAccount(t.namespace, kubernetes.NewServiceAccount(t.serviceAccountName))
		Expect(err).ToNot(HaveOccurred())

		t.podSecurityPolicyName = "deployment-psp"
		_, err = t.k8sClient.CreatePodSecurityPolicy(kubernetes.NewPodSecurityPolicy(t.podSecurityPolicyName, kubernetes.NewPodSecurityPolicySpec()))
		// tolerate pre-existing psp in order to make re-testing on previous environment idempotent
		if err != nil && err.Error() != "podsecuritypolicies.policy \"deployment-psp\" already exists" {
			Fail("got an error trying to create security policy: " + err.Error())
		}

		t.roleName = "deployment-role"
		_, err = t.k8sClient.CreateRole(t.namespace, kubernetes.NewRole(t.roleName, t.podSecurityPolicyName))
		Expect(err).ToNot(HaveOccurred())

		t.roleBindingName = "deployment-role-binding"
		_, err = t.k8sClient.CreateRoleBinding(t.namespace, kubernetes.NewRoleBinding(t.roleName, t.roleBindingName, t.serviceAccountName))
		Expect(err).ToNot(HaveOccurred())
	})

	By("Deploying workload 1 and 2", func() {
		var err error
		t.nginx1Deployment = kubernetes.NewDeployment("nginx-1", kubernetes.NewNginxDeploymentSpec(t.serviceAccountName))
		t.nginx1Deployment, err = t.k8sClient.CreateDeployment(t.namespace, t.nginx1Deployment)
		Expect(err).ToNot(HaveOccurred())

		t.nginx2Deployment = kubernetes.NewDeployment("nginx-2", kubernetes.NewNginxDeploymentSpec(t.serviceAccountName))
		t.nginx2Deployment, err = t.k8sClient.CreateDeployment(t.namespace, t.nginx2Deployment)
		Expect(err).ToNot(HaveOccurred())

		err = t.k8sClient.WaitForDeployment(t.namespace, t.nginx1Deployment.Name, t.timeout, GinkgoWriter)
		Expect(err).NotTo(HaveOccurred())

		err = t.k8sClient.WaitForDeployment(t.namespace, t.nginx2Deployment.Name, t.timeout, GinkgoWriter)
		Expect(err).NotTo(HaveOccurred())
	})
}

func (t *Deployment) AfterBackup(config Config) {
	By("Deleting workload 2", func() {
		err := t.k8sClient.DeleteDeployment(t.namespace, t.nginx2Deployment.Name)
		Expect(err).ToNot(HaveOccurred())
	})

	By("Deploying workload 3", func() {
		var err error
		t.nginx3Deployment = kubernetes.NewDeployment("nginx-3", kubernetes.NewNginxDeploymentSpec(t.serviceAccountName))
		t.nginx3Deployment, err = t.k8sClient.CreateDeployment(t.namespace, t.nginx3Deployment)
		Expect(err).ToNot(HaveOccurred())

		err = t.k8sClient.WaitForDeployment(t.namespace, t.nginx3Deployment.Name, t.timeout, GinkgoWriter)
		Expect(err).NotTo(HaveOccurred())
	})
}

func (t *Deployment) AfterRestore(config Config) {
	By("Waiting for workloads 1 and 2 to be available", func() {
		err := t.k8sClient.WaitForDeployment(t.namespace, t.nginx1Deployment.Name, t.timeout, GinkgoWriter)
		Expect(err).NotTo(HaveOccurred())

		err = t.k8sClient.WaitForDeployment(t.namespace, t.nginx2Deployment.Name, t.timeout, GinkgoWriter)
		Expect(err).NotTo(HaveOccurred())
	})

	By("Asserting that workload 3 is gone", func() {
		_, err := t.k8sClient.GetDeployment(t.namespace, t.nginx3Deployment.Name)
		Expect(err).To(HaveOccurred())

		statusErr, ok := err.(*errors.StatusError)
		Expect(ok).To(BeTrue())
		Expect(statusErr.ErrStatus.Code).To(Equal(int32(404)))
	})
}

func (t *Deployment) Cleanup(config Config) {
	By("Deleting the test namespace", func() {
		err := t.k8sClient.DeleteNamespace(t.namespace)
		Expect(err).NotTo(HaveOccurred())
	})

	By("Deleting the pod security policy", func() {
		err := t.k8sClient.DeletePodSecurityPolicy(t.podSecurityPolicyName)
		Expect(err).NotTo(HaveOccurred())
	})
}
