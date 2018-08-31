package acceptance

import (
	"testing"

	"bytes"
	"encoding/json"
	"fmt"
	"github.com/cloudfoundry-incubator/kubo-disaster-recovery-acceptance-tests/testcases"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"
	"io/ioutil"
	"os"
	"os/exec"
	"time"
)

var (
	config         Config
	artifactPath   string
	testCaseConfig = testcases.Config{}
	testCases      []TestCase
)

func TestAcceptance(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Acceptance Suite")
}

var _ = BeforeSuite(func() {
	mustHaveEnv("BOSH_ENVIRONMENT")
	mustHaveEnv("BOSH_CLIENT")
	mustHaveEnv("BOSH_CLIENT_SECRET")
	mustHaveEnv("BOSH_CA_CERT")
	mustHaveEnv("BOSH_DEPLOYMENT")

	ensureBBR()

	config = parseConfig(mustHaveEnv("CONFIG_PATH"))

	var err error
	artifactPath, err = ioutil.TempDir("", "k-drats")
	Expect(err).NotTo(HaveOccurred())

	SetDefaultEventuallyTimeout(config.TimeoutMinutes * time.Minute)

	testCases = []TestCase{
		testcases.FakeTestCase{},
	}

	fmt.Println("Running testcases: ")
	for _, t := range testCases {
		fmt.Println(t.Name())
	}
})

var _ = AfterSuite(func() {
	err := os.RemoveAll(artifactPath)
	Expect(err).NotTo(HaveOccurred())
})

func ensureBBR() {
	cmd := exec.Command("bbr")
	session, err := gexec.Start(cmd, new(bytes.Buffer), GinkgoWriter)
	Expect(err).NotTo(HaveOccurred())
	Eventually(session).Should(gexec.Exit(0))
}

func parseConfig(path string) Config {
	rawConfig, err := ioutil.ReadFile(path)
	Expect(err).NotTo(HaveOccurred())

	var config Config
	err = json.Unmarshal(rawConfig, &config)
	Expect(err).NotTo(HaveOccurred())
	return config
}

func mustHaveEnv(name string) string {
	val := os.Getenv(name)
	Expect(val).NotTo(BeEmpty(), fmt.Sprintf("Env var '%s' not set", name))
	return val
}
