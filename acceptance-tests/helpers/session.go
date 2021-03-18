package helpers

import (
	"fmt"
	. "github.com/onsi/gomega/gexec"
)

func SessionPrinter(message string, session *Session) func()string {
	return func() string {
		return fmt.Sprintf(
			"%s\nCMD: %s\nSTDOUT: %s\nSTDERR: %s\n",
			message,
			session.Command,
			string(session.Out.Contents()),
			string(session.Err.Contents()),
		)
	}
}

