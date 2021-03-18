package helpers

import (
	"acceptancetests/helpers/shell"

	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gexec"
)

func RestartApp(appName string) {
	session := shell.StartCommand(cf, "restart", appName)
	Eventually(session).Should(Exit(0), func() string {
		shell.RunCommand(cf, "logs", appName, "--recent")
		return "failed to restart app!"
	})
}
