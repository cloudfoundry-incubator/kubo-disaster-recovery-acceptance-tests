package acceptance

import (
	"testing"

	"github.com/cloudfoundry-incubator/kubo-disaster-recovery-acceptance-tests/command"

	"github.com/cloudfoundry-incubator/kubo-disaster-recovery-acceptance-tests/testcase"

	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"
)

const clusterName = "k-drats"

var (
	artifactPath   string
	kubeCACertPath string
	testCaseConfig = testcase.Config{}
	testCases      []TestCase
	filter         ConfigTestCaseFilter
)

func TestAcceptance(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Acceptance Suite")
}

var _ = BeforeSuite(func() {
	mustHaveBoshEnvVars()
	ensureBBR()

	config := NewConfig(mustHaveEnv("CONFIG_PATH"))
	filter = NewConfigTestCaseFilter(mustHaveEnv("CONFIG_PATH"))

	SetDefaultEventuallyTimeout(config.TimeoutMinutes * time.Minute)

	artifactPath = createTempDir()
	kubeCACertPath = writeTempFile(config.CACert)
	configureKubectl(config, kubeCACertPath)

	testCases = []TestCase{
		testcase.Deployment{},
		testcase.EtcdCluster{},
	}

	fmt.Println("Running test cases:")
	for _, t := range filter.Filter(testCases) {
		fmt.Println(t.Name())
	}
	fmt.Println("")
})

var _ = AfterSuite(func() {
	err := os.RemoveAll(kubeCACertPath)
	Expect(err).NotTo(HaveOccurred())

	err = os.RemoveAll(artifactPath)
	Expect(err).NotTo(HaveOccurred())
})

func mustHaveBoshEnvVars() {
	mustHaveEnv("BOSH_ENVIRONMENT")
	mustHaveEnv("BOSH_CLIENT")
	mustHaveEnv("BOSH_CLIENT_SECRET")
	mustHaveEnv("BOSH_CA_CERT")
	mustHaveEnv("BOSH_DEPLOYMENT")
}

func ensureBBR() {
	cmd := exec.Command("bbr")
	session, err := gexec.Start(cmd, new(bytes.Buffer), GinkgoWriter)
	Expect(err).NotTo(HaveOccurred())
	Eventually(session).Should(gexec.Exit(0))
}

func mustHaveEnv(name string) string {
	val := os.Getenv(name)
	Expect(val).NotTo(BeEmpty(), fmt.Sprintf("Env var '%s' not set", name))
	return val
}

func getArtifactFromPath(artifactPath string) string {
	files, err := ioutil.ReadDir(artifactPath)
	Expect(err).NotTo(HaveOccurred())
	Expect(files).To(HaveLen(1))

	return files[0].Name()
}

func writeTempFile(content string) string {
	tempFile, err := ioutil.TempFile("", "k-drats")
	Expect(err).NotTo(HaveOccurred())

	_, err = tempFile.Write([]byte(content))
	Expect(err).NotTo(HaveOccurred())

	return tempFile.Name()
}

func createTempDir() string {
	path, err := ioutil.TempDir("", "k-drats")
	Expect(err).NotTo(HaveOccurred())

	return path
}

func configureKubectl(config Config, caCertPath string) {
	command.RunSuccessfully(
		"kubectl config set-cluster",
		"kubectl",
		"config",
		"set-cluster", clusterName,
		fmt.Sprintf("--server=%s", config.APIServerURL),
		fmt.Sprintf("--certificate-authority=%s", caCertPath),
		"--embed-certs=true",
	)

	command.RunSuccessfully(
		"kubectl config set-credentials",
		"kubectl",
		"config",
		"set-credentials",
		config.Username,
		fmt.Sprintf("--token=%s", config.Password),
	)

	command.RunSuccessfully(
		"kubectl config set-context",
		"kubectl",
		"config",
		"set-context", clusterName,
		fmt.Sprintf("--cluster=%s", clusterName),
		fmt.Sprintf("--user=%s", config.Username),
	)

	command.RunSuccessfully(
		"kubectl config use-context",
		"kubectl",
		"config",
		"use-context", clusterName,
	)
}
