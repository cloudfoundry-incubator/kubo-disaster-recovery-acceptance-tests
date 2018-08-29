package acceptance_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/cloudfoundry-incubator/kubo-disaster-recovery-acceptance-tests/runner"
	"github.com/cloudfoundry-incubator/kubo-disaster-recovery-acceptance-tests/testcases"
)

var _ = Describe("Kubo", func() {
	config, err := runner.NewConfig()
	if err != nil {
		panic(err)
	}

	SetDefaultEventuallyTimeout(config.Timeout)

	testCases := []runner.TestCase{
		testcases.FakeTestCase{},
	}

	runner.RunBoshDisasterRecoveryAcceptanceTests(config, testCases)
})
