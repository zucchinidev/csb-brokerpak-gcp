package acceptance_test

import (
	"testing"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestAcceptance(t *testing.T) {
	RegisterFailHandler(Fail)
	SetDefaultEventuallyTimeout(10 * time.Minute)
	SetDefaultEventuallyPollingInterval(6 * time.Second)
	RunSpecs(t, "Acceptance Tests Suite")
}
