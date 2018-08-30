package acceptance

import "github.com/cloudfoundry-incubator/kubo-disaster-recovery-acceptance-tests/testcases"

type TestCase interface {
	Name() string
	BeforeBackup(testcases.Config)
	AfterBackup(testcases.Config)
	AfterRestore(testcases.Config)
	Cleanup(testcases.Config)
}
