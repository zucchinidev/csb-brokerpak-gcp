package main_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"net"
	"testing"
)

func TestTerraformProviderCSBPG(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "TerraformProviderCSBPG Suite")
}

var _ = BeforeSuite(func() {

})

var _ = AfterSuite(func() {

})

func freePort() int {
	listener, err := net.Listen("tcp", "localhost:0")
	Expect(err).NotTo(HaveOccurred())
	defer listener.Close()
	return listener.Addr().(*net.TCPAddr).Port
}
