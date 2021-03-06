package command

import (
	"fmt"
	"io"
	"io/ioutil"
	"os/exec"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"
)

func RunSuccessfully(description, cmd string, args ...string) *gexec.Session {
	session := runCommand(description, GinkgoWriter, cmd, args...)
	Expect(session).To(gexec.Exit(0), "Command errored: "+description)
	return session
}

func RunSuccessfullyWithoutStream(description, cmd string, args ...string) *gexec.Session {
	session := runCommand(description, ioutil.Discard, cmd, args...)
	Expect(session).To(gexec.Exit(0), "Command errored: "+description)
	return session
}

func runCommand(description string, writer io.Writer, cmd string, args ...string) *gexec.Session {
	command := exec.Command(cmd, args...)
	session, err := gexec.Start(command, writer, writer)
	Expect(err).NotTo(HaveOccurred())

	Eventually(session).Should(gexec.Exit(), "Command timed out: "+description)
	fmt.Fprintln(writer, "")
	return session
}
