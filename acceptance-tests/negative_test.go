package acceptance_test

//var _ = Describe("Negative Responses", func() {
//	XDescribeTable(
//		"create service failures",
//		func(offering, plan string) {
//			name := helpers.RandomString()
//			defer helpers.DeleteService(name)
//
//			session := helpers.StartCreateService(offering, plan, name, `{"region":"bogus"}`)
//			session.Wait()
//			Expect(session).NotTo(Exit(0), helpers.SessionPrinter("expected command to fail, but it unexpectedly succeeded!", session))
//		},
//		Entry("redis", "csb-google-redis", "basic"),
//		Entry("mysql", "csb-google-mysql", "medium"),
//	)
//})
