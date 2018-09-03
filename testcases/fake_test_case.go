package testcases

import "github.com/cloudfoundry-incubator/kubo-disaster-recovery-acceptance-tests/helpers"

type FakeTestCase struct{}

func (t FakeTestCase) Name() string {
	return "fake_test_case"
}

func (t FakeTestCase) BeforeBackup(config Config) {
	helpers.RunCommandSuccessfullyWithFailureMessage("kubectl get all", "kubectl get all")
}

func (t FakeTestCase) AfterBackup(config Config) {}

func (t FakeTestCase) AfterRestore(config Config) {}

func (t FakeTestCase) Cleanup(config Config) {}
