package acceptance

import (
	"testing"
	"time"

	"github.com/cloudfoundry-incubator/kubo-disaster-recovery-acceptance-tests/testcase"

	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"
)

var (
	artifactPath   string
	kubeCACertPath string
	testCaseConfig = testcase.Config{}
	testCases      []TestCase
	filter         TestCaseFilter
)

var availableTestCases = []TestCase{
	testcase.NewDeployment(),
	testcase.NewEtcdCluster(),
	testcase.NewPodHealth(),
}

func TestAcceptance(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Acceptance Suite")
}

var _ = BeforeSuite(func() {
	mustHaveBoshEnvVars()
	ensureBBR()

	testCases = availableTestCases
	config := NewConfig(os.Getenv("CONFIG_PATH"))
	filter = NewTestCaseFilter(os.Getenv("CONFIG_PATH"))
	if filter != nil {
		testCases = filter.Filter(availableTestCases)
	}

	fmt.Println("Running test cases:")
	for _, t := range testCases {
		fmt.Println("* ", t.Name())
	}
	fmt.Println("")

	SetDefaultEventuallyTimeout(time.Minute * config.TimeoutMinutes)
	SetDefaultEventuallyPollingInterval(time.Second * 5)
	fmt.Printf("Timeout: %d min\n\n", config.TimeoutMinutes)
	artifactPath = createTempDir()
})

var _ = AfterSuite(func() {
	err := os.RemoveAll(artifactPath)
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

func createTempDir() string {
	path, err := ioutil.TempDir("", "k-drats")
	Expect(err).NotTo(HaveOccurred())

	return path
}
