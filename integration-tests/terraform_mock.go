package integration_tests

import (
	"encoding/json"
	"fmt"
	"github.com/onsi/gomega/gexec"
	"io/ioutil"
	"os"
	"path"
	"strings"
)

type TerraformMock struct {
	binary          string
	invocationStore string
}

func (p TerraformMock) ApplyInvocations() ([]TerraformInvocation, error) {
	invocations, err := p.Invocations()
	if err != nil {
		return nil, err
	}
	filteredInovocations := []TerraformInvocation{}
	for _, invocation := range invocations {
		if invocation.Type == "apply" {
			filteredInovocations = append(filteredInovocations, invocation)
		}
	}
	return filteredInovocations, nil
}

type TerraformInvocation struct {
	Type string
	dir  string
}

func (i TerraformInvocation) TFVars() (map[string]interface{}, error) {
	output := map[string]interface{}{}
	file, err := os.Open(path.Join(i.dir, "terraform.tfvars.json"))
	if err != nil {
		return nil, err
	}

	return output, json.NewDecoder(file).Decode(&output)
}

func (p TerraformMock) Invocations() ([]TerraformInvocation, error) {
	fileInfo, err := ioutil.ReadDir(p.invocationStore)
	if err != nil {
		return nil, err
	}
	invocations := []TerraformInvocation{}

	for _, file := range fileInfo {
		parts := strings.Split(file.Name(), "-")
		invocations = append(invocations, TerraformInvocation{Type: parts[0], dir: path.Join(p.invocationStore, file.Name())})
	}
	return invocations, nil
}

func (p TerraformMock) Reset() error {
	dir, err := ioutil.ReadDir(p.invocationStore)
	if err != nil {
		return err
	}
	for _, d := range dir {
		if err := os.RemoveAll(path.Join(p.invocationStore, d.Name())); err != nil {
			return err
		}
	}
	return nil
}

func NewTerraformMock() (TerraformMock, error) {
	dir, err := ioutil.TempDir("", "invocation_store")
	if err != nil {
		return TerraformMock{}, err
	}
	build, err := gexec.Build("github.com/cloudfoundry-incubator/csb-brokerpak-gcp/integration-tests/testframework/mock-binary/terraform", "-ldflags", fmt.Sprintf("-X 'main.InvocationStore=%s'", dir))
	if err != nil {
		return TerraformMock{}, err
	}
	return TerraformMock{binary: build, invocationStore: dir}, nil
}
