package integration_tests

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Mysql", Ordered, func() {
	var mockTerraform TerraformMock
	var broker TestCSBInstance

	BeforeAll(func() {
		var err error
		mockTerraform, err = NewTerraformMock()
		Expect(err).NotTo(HaveOccurred())
		broker = BuildBrokerWithProvider(mockTerraform)
		broker.Start()
	})
	AfterEach(func() {
		Expect(mockTerraform.Reset()).NotTo(HaveOccurred())
	})
	It("publish mysql in the catalog", func() {
		catalog, err := broker.Catalog()
		Expect(err).NotTo(HaveOccurred())
		service := FindService(catalog, "csb-google-mysql")
		Expect(service.Plans).To(HaveLen(3))
		Expect(service.Tags).To(ContainElement("preview"))
		Expect(service.Metadata.ImageUrl).NotTo(BeNil())
	})

	It("should provision small plan", func() {
		broker.Provision("csb-google-mysql", "small", nil)

		invocations, err := mockTerraform.ApplyInvocations()
		Expect(err).NotTo(HaveOccurred())
		Expect(invocations).To(HaveLen(1))

		Expect(invocations[0].TFVars()).To(HaveKeyWithValue("db_name", "csb-db"))
		Expect(invocations[0].TFVars()).To(HaveKeyWithValue("database_version", "MYSQL_5_7"))
		Expect(invocations[0].TFVars()).To(HaveKeyWithValue("cores", float64(2)))
		Expect(invocations[0].TFVars()).To(HaveKeyWithValue("storage_gb", float64(10)))

	})

	It("user should not allow changing mysql cores", func() {
		err := broker.Provision("csb-google-mysql", "small", map[string]interface{}{"cores": 5})
		Expect(err).To(HaveOccurred())
		Expect(err).To(MatchError(ContainSubstring("plan defined properties cannot be changed: cores")))
	})

	It("user should be able to update database name", func() {
		err := broker.Provision("csb-google-mysql", "small", map[string]interface{}{"db_name": "foobar"})

		invocations, err := mockTerraform.ApplyInvocations()
		Expect(err).NotTo(HaveOccurred())
		Expect(invocations).To(HaveLen(1))

		Expect(invocations[0].TFVars()).To(HaveKeyWithValue("db_name", "foobar"))
	})
})
