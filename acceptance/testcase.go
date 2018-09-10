package acceptance

import (
	"github.com/cloudfoundry-incubator/kubo-disaster-recovery-acceptance-tests/testcase"
)

type TestCase interface {
	Name() string
	BeforeBackup(testcase.Config)
	AfterBackup(testcase.Config)
	AfterRestore(testcase.Config)
	Cleanup(testcase.Config)
}
