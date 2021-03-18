package helpers

import (
	"fmt"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gexec"
	"os/exec"
)

func Run(path string, args ...string) {
	command := exec.Command(path, args...)
	session, err := Start(command, GinkgoWriter, GinkgoWriter)
	Expect(err).NotTo(HaveOccurred(), fmt.Sprintf("Error starting command: %s\n", command))
	Eventually(session).Should(Exit(0),SessionPrinter("error running command", session))
}
