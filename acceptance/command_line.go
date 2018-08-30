package acceptance

import (
	"fmt"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"
	"os/exec"
)

func runCommandSuccessfullyWithFailureMessage(description string, cmd string) *gexec.Session {
	session := runCommandWithStream(description, cmd)
	Expect(session).To(gexec.Exit(0), "Command errored: "+description)
	return session
}

func runCommandWithStream(description string, cmd string) *gexec.Session {
	command := exec.Command("bash", "-c", cmd)
	session, err := gexec.Start(command, GinkgoWriter, GinkgoWriter)
	Expect(err).ToNot(HaveOccurred())

	Eventually(session).Should(gexec.Exit(), "Command timed out: "+description)
	fmt.Fprintln(GinkgoWriter, "")
	return session
}
