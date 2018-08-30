package acceptance

import (
	"testing"

	"fmt"
	"github.com/cloudfoundry-incubator/kubo-disaster-recovery-acceptance-tests/testcases"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"time"
)

var config Config
var testCaseConfig testcases.Config
var testCases []TestCase

func TestAcceptance(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Acceptance Suite")
}

var _ = BeforeSuite(func() {
	var err error
	config, err = NewConfig()
	if err != nil {
		panic(err)
	}

	testCaseConfig = testcases.Config{}
	SetDefaultEventuallyTimeout(config.TimeoutMinutes * time.Minute)

	testCases = []TestCase{
		testcases.FakeTestCase{},
	}

	fmt.Println("Running testcases: ")
	for _, t := range testCases {
		fmt.Println(t.Name())
	}
})
