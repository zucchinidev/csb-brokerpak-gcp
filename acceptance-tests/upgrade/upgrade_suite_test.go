package upgrade_test

import (
	"csbbrokerpakgcp/acceptance-tests/helpers/environment"
	"flag"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var (
	developmentBuildDir string
	releasedBuildDir    string
	metadata            environment.GCPMetadata
)

func init() {
	flag.StringVar(&releasedBuildDir, "releasedBuildDir", "../../../gcp-released", "location of released version of built broker and brokerpak")
	flag.StringVar(&developmentBuildDir, "developmentBuildDir", "../../dev-release", "location of development version of built broker and brokerpak")
}

var _ = BeforeSuite(func() {
	metadata = environment.ReadGCPMetadata()
})

func TestUpgrade(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Upgrade Suite")
}
