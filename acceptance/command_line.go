package acceptance

import (
	"fmt"
	"os/exec"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"
)

func RunCommandSuccessfullyWithFailureMessage(description string, cmd string, args ...string) *gexec.Session {
	session := runCommandWithStream(description, cmd, args...)
	Expect(session).To(gexec.Exit(0), "Command errored: "+description)
	return session
}

func runCommandWithStream(description string, cmd string, args ...string) *gexec.Session {
	command := exec.Command(cmd, args...)
	session, err := gexec.Start(command, GinkgoWriter, GinkgoWriter)
	Expect(err).ToNot(HaveOccurred())

	Eventually(session).Should(gexec.Exit(), "Command timed out: "+description)
	fmt.Fprintln(GinkgoWriter, "")
	return session
}