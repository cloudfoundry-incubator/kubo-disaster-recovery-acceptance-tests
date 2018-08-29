package testcases

import "github.com/cloudfoundry-incubator/kubo-disaster-recovery-acceptance-tests/runner"

type FakeTestCase struct{}

func (t FakeTestCase) Name() string {
	return "fake_test_case"
}

func (t FakeTestCase) BeforeBackup(config runner.Config) {}

func (t FakeTestCase) AfterBackup(config runner.Config) {}

func (t FakeTestCase) AfterRestore(config runner.Config) {}

func (t FakeTestCase) Cleanup(config runner.Config) {}
