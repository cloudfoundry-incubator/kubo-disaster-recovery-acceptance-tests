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
	"github.com/cloudfoundry-incubator/kubo-disaster-recovery-acceptance-tests/helpers"
)

var (
	config         Config
	artifactPath   string
	testCaseConfig = testcases.Config{}
	testCases      []TestCase
	kubeCACertPath string
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

	setKubectlConfig(config)

	var err error
	artifactPath, err = ioutil.TempDir("", "k-drats")
	Expect(err).NotTo(HaveOccurred())

	SetDefaultEventuallyTimeout(config.TimeoutMinutes * time.Minute)

	testCases = []TestCase{
		testcases.FakeTestCase{},
		testcases.KuboTestCase{},
	}

	fmt.Println("Running testcases: ")
	for _, t := range testCases {
		fmt.Println(t.Name())
	}
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
	kubeCACertFile, err := ioutil.TempFile("","kubeCACert")
	Expect(err).NotTo(HaveOccurred())

	_, err = kubeCACertFile.Write([]byte(config.CACert))
	Expect(err).NotTo(HaveOccurred())
	kubeCACertPath = kubeCACertFile.Name()

	helpers.RunCommandSuccessfullyWithFailureMessage(
		"kubectl config set-cluster",
		fmt.Sprintf(
			"kubectl config set-cluster %s --server=%s --certificate-authority=%s --embed-certs=true",
			config.ClusterName,
			config.APIServerURL,
			kubeCACertPath,
		),
	)
	helpers.RunCommandSuccessfullyWithFailureMessage(
		"kubectl config set-credentials",
		fmt.Sprintf(
			"kubectl config set-credentials %s --token='%s'",
			config.Username,
			config.Password,
		),
	)
	helpers.RunCommandSuccessfullyWithFailureMessage(
		"kubectl config set-context",
		fmt.Sprintf(
			"kubectl config set-context %s --cluster=%s --user=%s",
			config.ClusterName,
			config.ClusterName,
			config.Username,
		),
	)
	helpers.RunCommandSuccessfullyWithFailureMessage(
		"kubectl config use-context",
		fmt.Sprintf(
			"kubectl config use-context %s",
			config.ClusterName,
		),
	)
}
