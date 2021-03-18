package acceptance_test

import (
	"acceptancetests/helpers"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

const appName = "spring-music"

var _ = Describe("Redis", func() {
	var serviceInstanceName string

	BeforeEach(func() {
		serviceInstanceName = helpers.RandomString()
		helpers.CreateService("csb-google-redis", "basic", serviceInstanceName)
	})

	It("validates", func() {
		binding := helpers.Bind(appName, serviceInstanceName)

		By("checking that the credential is a credhub reference")
		creds := helpers.GetBindingCredential(appName, "csb-google-redis", binding)
		Expect(creds).To(HaveLen(1))
		Expect(creds).To(HaveKey("credhub-ref"))

		By("restarting spring music")
		helpers.RestartApp(appName)

		By("validating with Spring Music")
		app := helpers.PushSpringMusicValidator()
		defer helpers.DeleteApp(app)

		helpers.Bind(app, serviceInstanceName)
		helpers.RestartApp(app)
	})

	AfterEach(func() {
		helpers.Unbind(appName, serviceInstanceName)
		helpers.DeleteService(serviceInstanceName)
	})
})
