package helpers

import (
	"fmt"
	"os/exec"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gexec"
)

const cf = "cf"

func CF(args ...string) *Session {
	session := StartCF(args...)
	Eventually(session).Should(Exit(0), SessionPrinter("cf command failed", session))
	return session
}

func StartCF(args ...string) *Session {
	cf, err := exec.LookPath("cf")
	Expect(err).NotTo(HaveOccurred(), "cannot find cf command")

	command := exec.Command(cf, args...)
	// fmt.Sprintf(GinkgoWriter, command) ???
	session, err := Start(command, GinkgoWriter, GinkgoWriter)
	Expect(err).NotTo(HaveOccurred(), fmt.Sprintf("error starting: %s", command))
	return session
}
