package helpers

import (
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gbytes"
)

func DeleteService(name string) {
	CF("delete-service", "-f", name)
	Eventually(func() *Buffer { return CF("services").Out }).ShouldNot(Say(name))
}
