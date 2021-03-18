package helpers

import (
	"fmt"
	"strings"

	. "github.com/onsi/ginkgo"

	"code.cloudfoundry.org/jsonry"

	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gexec"
)

func GetBindingCredential(appName, serviceName, bindingName string) interface{} {
	session := CF("app", "--guid", appName)
	Eventually(session).Should(Exit(0))
	guid := strings.TrimSpace(string(session.Out.Contents()))

	session = CF("curl", fmt.Sprintf("/v3/apps/%s/env", guid))
	Eventually(session).Should(Exit(0))
	env := strings.TrimSpace(string(session.Out.Contents()))

	var receiver struct {
		Services map[string]interface{} `jsonry:"system_env_json.VCAP_SERVICES"`
	}
	err := jsonry.Unmarshal([]byte(env), &receiver)
	Expect(err).NotTo(HaveOccurred())

	Expect(receiver.Services).NotTo(BeEmpty())
	Expect(receiver.Services).To(HaveKey(serviceName))
	bindings := receiver.Services[serviceName]
	Expect(bindings).To(BeAssignableToTypeOf([]interface{}{}))
	Expect(bindings).NotTo(BeEmpty())

	for _, b := range bindings.([]interface{}) {
		if n, ok := b.(map[string]interface{})["name"]; ok && n == bindingName {
			Expect(b).To(HaveKey("credentials"))
			return b.(map[string]interface{})["credentials"]
		}
	}

	Fail(fmt.Sprintf("could not find data for binding: %s\n%+v", bindingName, bindings))
	return nil
}
