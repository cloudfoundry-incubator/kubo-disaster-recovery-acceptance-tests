package acceptance

import (
	"github.com/cloudfoundry-incubator/kubo-disaster-recovery-acceptance-tests/testcase"
	"testing"

	"bytes"
	"encoding/json"
	"fmt"
	"github.com/cloudfoundry-incubator/kubo-disaster-recovery-acceptance-tests/helpers"
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
	testCaseConfig = testcase.Config{}
	testCases      []TestCase
	kubeCACertPath string
	filter         ConfigTestCaseFilter
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
	filter = NewConfigTestCaseFilter(mustHaveEnv("CONFIG_PATH"))

	setKubectlConfig(config)

	var err error
	artifactPath, err = ioutil.TempDir("", "k-drats")
	Expect(err).NotTo(HaveOccurred())

	SetDefaultEventuallyTimeout(config.TimeoutMinutes * time.Minute)

	testCases = []TestCase{
		testcase.Deployment{},
	}

	fmt.Println("Running test cases: ")
	for _, t := range filter.Filter(testCases) {
		fmt.Println(t.Name())
	}
	fmt.Println("")
})

var _ = AfterSuite(func() {
	err := os.RemoveAll(artifactPath)
	Expect(err).NotTo(HaveOccurred())
	err = os.RemoveAll(kubeCACertPath)
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

func setKubectlConfig(config Config) {
	kubeCACertFile, err := ioutil.TempFile("", "kubeCACert")
	Expect(err).NotTo(HaveOccurred())

	_, err = kubeCACertFile.Write([]byte(config.CACert))
	Expect(err).NotTo(HaveOccurred())
	kubeCACertPath = kubeCACertFile.Name()

	RunCommandSuccessfullyWithFailureMessage(
		"kubectl config set-cluster",
		"kubectl", "config", "set-cluster", config.ClusterName,
		fmt.Sprintf("--server=%s", config.APIServerURL),
		fmt.Sprintf("--certificate-authority=%s", kubeCACertPath),
		"--embed-certs=true",
	)

	RunCommandSuccessfullyWithFailureMessage(
		"kubectl config set-credentials",
		"kubectl", "config", "set-credentials", config.Username, fmt.Sprintf("--token=%s", config.Password),
	)

	RunCommandSuccessfullyWithFailureMessage(
		"kubectl config set-context",
		"kubectl", "config", "set-context", config.ClusterName,
		fmt.Sprintf("--cluster=%s", config.ClusterName),
		fmt.Sprintf("--user=%s", config.Username),
	)

	RunCommandSuccessfullyWithFailureMessage(
		"kubectl config use-context",
		"kubectl", "config", "use-context", config.ClusterName,
	)
}
