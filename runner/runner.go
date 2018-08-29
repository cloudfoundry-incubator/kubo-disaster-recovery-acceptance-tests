package runner

import (
	"fmt"

	. "github.com/onsi/ginkgo"
)

func RunBoshDisasterRecoveryAcceptanceTests(config Config, testCases []TestCase) {
	fmt.Println("Running testcases: ")
	for _, t := range testCases {
		fmt.Println(t.Name())
	}

	It("can backup and restore", func() {
		By("running the before backup step", func() {
			for _, testCase := range testCases {
				fmt.Println("Running the before backup step for " + testCase.Name())
				testCase.BeforeBackup(config)
			}
		})

		By("backing up", func() {
			fmt.Println("bbr deployment -d cfcr backup")
		})

		By("running the after backup step", func() {
			for _, testCase := range testCases {
				fmt.Println("Running the after backup step for " + testCase.Name())
				testCase.AfterBackup(config)
			}
		})

		By("restoring", func() {
			fmt.Println("bbr deployment -d cfcr restore --artifact-path blah")
		})

		By("waiting for kubo api to be available", func() {
			fmt.Println("wait for kubo api")
		})

		By("running the after restore step", func() {
			for _, testCase := range testCases {
				fmt.Println("Running the after restore step for " + testCase.Name())
				testCase.AfterRestore(config)
			}
		})
	})

	AfterEach(func() {
		By("running bbr deployment backup-cleanup", func() {
			fmt.Println("bbr deployment -d cfcr backup-cleanup")
		})

		By("Running cleanup for each testcase", func() {
			for _, testCase := range testCases {
				fmt.Println("Running the cleanup step for " + testCase.Name())
				testCase.Cleanup(config)
			}
		})
	})
}
