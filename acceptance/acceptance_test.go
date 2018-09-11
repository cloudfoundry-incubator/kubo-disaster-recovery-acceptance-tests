package acceptance

import (
	"fmt"

	"github.com/cloudfoundry-incubator/kubo-disaster-recovery-acceptance-tests/kubernetes"

	"github.com/cloudfoundry-incubator/kubo-disaster-recovery-acceptance-tests/command"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Kubo", func() {
	It("can backup and restore", func() {
		By("running the before backup step", func() {
			for _, testCase := range testCases {
				fmt.Println("Running the before backup step for " + testCase.Name())
				testCase.BeforeBackup(testCaseConfig)
			}
		})

		By("backing up", func() {
			command.RunSuccessfully(
				"bbr deployment backup",
				"bbr", "deployment", "backup", "--artifact-path", artifactPath,
			)
		})

		By("running the after backup step", func() {
			for _, testCase := range testCases {
				fmt.Println("Running the after backup step for " + testCase.Name())
				testCase.AfterBackup(testCaseConfig)
			}
		})

		By("restoring", func() {
			artifact := getArtifactFromPath(artifactPath)
			command.RunSuccessfully(
				"bbr deployment restore",
				"bbr", "deployment", "restore", "--artifact-path", fmt.Sprintf("%s/%s", artifactPath, artifact),
			)
		})

		k8sClient, err := kubernetes.NewKubeClient()
		Expect(err).NotTo(HaveOccurred())

		By("Waiting for API to be available", func() {
			Eventually(func() bool {
				return k8sClient.IsHealthy()
			}, "60s", "5s").Should(BeTrue())

		})

		By("Waiting for system workloads", func() {
			expectedSelector := []string{"kube-dns", "heapster", "kubernetes-dashboard", "influxdb"}

			for _, selector := range expectedSelector {
				deployments, err := k8sClient.GetDeployments("kube-system", "k8s-app="+selector)
				Expect(err).NotTo(HaveOccurred())

				Expect(deployments.Items).To(HaveLen(1), fmt.Sprintf("one %s deployment should exist, instead found: %#v", selector, deployments.Items))

				err = k8sClient.WaitForDeployment("kube-system", deployments.Items[0].Name, GinkgoWriter)
				Expect(err).NotTo(HaveOccurred())
			}
		})

		By("running the after restore step", func() {
			for _, testCase := range testCases {
				fmt.Println("Running the after restore step for " + testCase.Name())
				testCase.AfterRestore(testCaseConfig)
			}
		})
	})

	AfterEach(func() {
		By("running bbr deployment backup-cleanup", func() {
			command.RunSuccessfully(
				"bbr deployment backup-cleanup",
				"bbr", "deployment", "backup-cleanup",
			)
		})

		By("Running cleanup for each testcase", func() {
			for _, testCase := range testCases {
				fmt.Println("Running the cleanup step for " + testCase.Name())
				testCase.Cleanup(testCaseConfig)
			}
		})
	})
})
