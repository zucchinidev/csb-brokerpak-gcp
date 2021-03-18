package helpers

import (
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gbytes"
	. "github.com/onsi/gomega/gexec"
)

func CreateService(offering, plan, name string, parameters ...string) {
	session := StartCreateService(offering, plan, name, parameters...)
	Eventually(session).Should(Exit(0), SessionPrinter("failed to create service", session))
	Eventually(func() *Buffer { return CF("service", name).Out }).Should(Say(`status:\s+create succeeded`))
}

func StartCreateService(offering, plan, name string, parameters ...string) *Session {
	args := []string{"create-service", offering, plan, name}
	if len(parameters) > 0 {
		args = append(args, "-c", parameters[0])
	}

	session := StartCF(args...)
	return session
}
