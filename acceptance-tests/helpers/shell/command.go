package shell

import (
	"fmt"
	"os/exec"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gexec"
)

func StartCommand(cmd string, args ...string) *Session {
	path, err := exec.LookPath(cmd)
	Expect(err).NotTo(HaveOccurred(), "cannot find cf command")

	command := exec.Command(path, args...)
	fmt.Fprintf(GinkgoWriter, "RUNNING: %s\n", command)
	session, err := Start(command, GinkgoWriter, GinkgoWriter)
	Expect(err).NotTo(HaveOccurred(), fmt.Sprintf("error starting: %s", command))
	return session
}

func RunCommand(cmd string, args ...string) ([]byte, []byte) {
	session := StartCommand(cmd, args...)
	Eventually(session).Should(Exit(0))
	return session.Out.Contents(), session.Err.Contents()
}
