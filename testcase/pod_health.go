package testcase

import (
	"os/exec"
	"time"

	"github.com/onsi/gomega/gbytes"
	"github.com/onsi/gomega/gexec"

	"k8s.io/api/apps/v1"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/cloudfoundry-incubator/kubo-disaster-recovery-acceptance-tests/kubernetes"
)

type PodHealth struct {
	k8sClient  *kubernetes.Client
	namespace  string
	deployment *v1.Deployment
	timeout    time.Duration
}

func (PodHealth) Name() string {
	return "pod_health"
}

func NewPodHealth() *PodHealth {
	return &PodHealth{}
}

func (p *PodHealth) BeforeBackup(config Config) {
	By("Initializing K8s client", func() {
		var err error
		p.k8sClient, err = kubernetes.NewKubeClient()
		Expect(err).NotTo(HaveOccurred())

		nsObject, err := p.k8sClient.CreateNamespace("bbr")
		Expect(err).ToNot(HaveOccurred())

		p.namespace = nsObject.Name
		p.timeout = 5 * time.Minute
	})

	By("Deploying nginx", func() {
		var err error
		p.deployment = kubernetes.NewDeployment("nginx", kubernetes.NewNginxDeploymentSpec())
		p.deployment, err = p.k8sClient.CreateDeployment(p.namespace, p.deployment)
		Expect(err).ToNot(HaveOccurred())

		err = p.k8sClient.WaitForDeployment(p.namespace, p.deployment.Name, p.timeout, GinkgoWriter)
		Expect(err).NotTo(HaveOccurred())
	})
}

func (p *PodHealth) AfterBackup(config Config) {
	By("Deleting nginx", func() {
		err := p.k8sClient.DeleteDeployment(p.namespace, p.deployment.Name)
		Expect(err).ToNot(HaveOccurred())
	})
}

func (p *PodHealth) AfterRestore(config Config) {
	var podName string
	By("Waiting for pod to be present", func() {
		p.k8sClient.WaitForDeployment(p.namespace, p.deployment.Name, p.timeout, GinkgoWriter)
	})

	By("Allowing commands to be executed on the container", func() {
		args := []string{"get", "pods", "-l", "app=" + p.deployment.Name, "-o", "jsonpath={.items[0].metadata.name}"}
		session := runKubectlCommandInNamespace(p.namespace, args...)
		Eventually(session, "15s").Should(gexec.Exit(0))
		podName = string(session.Out.Contents())

		execArgs := []string{"exec", podName, "--", "which", "nginx"}
		execSession := runKubectlCommandInNamespace(p.namespace, execArgs...)
		Eventually(execSession, "60s").Should(gexec.Exit(0))
		Expect(execSession.Out).To(gbytes.Say("/usr/sbin/nginx"))
	})

	By("Allowing access to pod logs", func() {
		port := "57869"
		args := []string{"port-forward", podName, port + ":80"}
		portForwardSess := runKubectlCommandInNamespace(p.namespace, args...)

		Eventually(curlFunc("http://localhost:"+port), "15s").Should(ContainSubstring("Server: nginx"))
		getLogs := runKubectlCommandInNamespace(p.namespace, "logs", podName)
		Eventually(getLogs, "15s").Should(gexec.Exit(0))
		logContent := string(getLogs.Out.Contents())

		Expect(logContent).To(ContainSubstring("curl"))
		portForwardSess.Terminate().Wait("15s")
	})
}

func (p *PodHealth) Cleanup(config Config) {
	By("Deleting the test namespace and port forwarding", func() {
		err := p.k8sClient.DeleteNamespace(p.namespace)
		Expect(err).NotTo(HaveOccurred())
	})
}

func runKubectlCommandInNamespace(namespace string, args ...string) *gexec.Session {
	newArgs := append([]string{"--namespace", namespace}, args...)
	command := exec.Command("kubectl", newArgs...)
	session, err := gexec.Start(command, GinkgoWriter, GinkgoWriter)

	Expect(err).NotTo(HaveOccurred())
	return session
}

func curlFunc(endpoint string) func() (string, error) {
	return func() (string, error) {
		cmd := exec.Command("curl", "--head", endpoint)
		out, err := cmd.CombinedOutput()
		return string(out), err
	}
}
