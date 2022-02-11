package integration_tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	_ "github.com/cloudfoundry-incubator/cloud-service-broker/cmd"
	"github.com/google/uuid"
	"github.com/hashicorp/go-retryablehttp"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"
	"github.com/pivotal-cf/brokerapi/v8/domain"
	"github.com/pivotal-cf/brokerapi/v8/domain/apiresponses"
	"github.com/pkg/errors"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"text/template"
	"time"
)

type TestCSBInstance struct {
	brokerBuild string
	workspace   string
	password    string
	username    string
	port        string
}

func (instance TestCSBInstance) Start() error {
	file, err := ioutil.TempFile("", "test-db")
	if err != nil {
		return err
	}

	go func() {
		serverCommand := exec.Command(instance.brokerBuild)
		serverCommand.Dir = instance.workspace
		serverCommand.Env = []string{
			"DB_PATH=" + file.Name(),
			"DB_TYPE=sqlite3",
			"PORT=" + instance.port,
			"SECURITY_USER_NAME=" + instance.username,
			"SECURITY_USER_PASSWORD=" + instance.password,
			"GOOGLE_CREDENTIALS=credentials",
			"GOOGLE_PROJECT=project-name",
			"GCP_CREDENTIALS=credentials",
			"GCP_PROJECT=project",
		}
		serverCommand.Stdout = GinkgoWriter
		serverCommand.Stderr = GinkgoWriter
		err := serverCommand.Run()
		panic(errors.Wrap(err, fmt.Sprintf("failed starting broker on workspace %s, with build %s", instance.workspace, instance.brokerBuild)))
	}()

	return waitForHttpServer("http://localhost:" + instance.port)
}

func waitForHttpServer(s string) error {
	retryClient := retryablehttp.NewClient()
	retryClient.RetryMax = 10
	retryClient.RetryWaitMin = time.Second * 2
	retryClient.RetryWaitMax = time.Second * 5
	_, err := retryClient.Get(s)
	return err
}

func (instance TestCSBInstance) Provision(serviceName string, planName string, params map[string]interface{}) error {
	instanceID, resp, err := instance.provision(serviceName, planName, params)
	if err != nil {
		return err
	}

	return instance.pollLastOperation(instanceID, resp.OperationData)
}

func (instance TestCSBInstance) provision(serviceName string, planName string, params map[string]interface{}) (string, *apiresponses.ProvisioningResponse, error) {
	catalog, err := instance.Catalog()
	if err != nil {
		return "", nil, err
	}
	serviceGuid, planGuid := FindServicePlanGUIDs(catalog, serviceName, planName)
	details := domain.ProvisionDetails{
		ServiceID: serviceGuid,
		PlanID:    planGuid,
	}
	if params != nil {
		data, err := json.Marshal(&params)
		if err != nil {
			return "", nil, err
		}
		details.RawParameters = json.RawMessage(data)
	}
	data, err := json.Marshal(details)
	if err != nil {
		return "", nil, err
	}
	instanceId := uuid.New()
	body, status, err := instance.httpInvokeBroker("service_instances/"+instanceId.String()+"?accepts_incomplete=true", "PUT", bytes.NewBuffer(data))
	if err != nil {
		return "", nil, err
	}
	if status != http.StatusAccepted {
		return "", nil, fmt.Errorf("request failed: status %d: body %s", status, body)
	}

	response := apiresponses.ProvisioningResponse{}

	return instanceId.String(), &response, json.Unmarshal(body, &response)
}

func (instance TestCSBInstance) pollLastOperation(instanceID string, lastOperation string) error {
	timeout := time.After(5 * time.Second)
	ticker := time.NewTicker(500 * time.Millisecond)
	defer ticker.Stop()

	for true {
		select {
		case <-timeout:
			return fmt.Errorf("timed out polling %s %s", instanceID, lastOperation)
		case <-ticker.C:
			data, status, err := instance.httpInvokeBroker("service_instances/"+instanceID+"/last_operation", "GET", nil)
			if err != nil {
				return err
			}
			if status != http.StatusOK {
				return fmt.Errorf("request failed: status %d: body %s", status, data)
			}
			resp := apiresponses.LastOperationResponse{}
			err = json.Unmarshal(data, &resp)
			if err != nil {
				return err
			}
			if resp.State != domain.InProgress {
				return nil
			}
		}
	}
	return nil
}

func (instance TestCSBInstance) Catalog() (*apiresponses.CatalogResponse, error) {
	catalogJson, status, err := instance.httpInvokeBroker("catalog", "GET", nil)
	if err != nil {
		return nil, err
	}
	if status != http.StatusOK {
		return nil, fmt.Errorf("request failed: status %d: body %s", status, catalogJson)
	}

	resp := &apiresponses.CatalogResponse{}
	return resp, json.Unmarshal(catalogJson, resp)
}

func (instance TestCSBInstance) httpInvokeBroker(subpath string, method string, body io.Reader) ([]byte, int, error) {
	client := &http.Client{
		Timeout: time.Second * 10,
	}
	req, err := http.NewRequest(method, instance.BrokerUrl(subpath), body)
	if err != nil {
		return nil, 0, err
	}
	req.Header.Set("X-Broker-API-Version", "2.14")
	req.SetBasicAuth(instance.username, instance.password)
	response, err := client.Do(req)
	if err != nil {
		return nil, 0, err
	}
	defer response.Body.Close()
	contents, err := io.ReadAll(response.Body)
	return contents, response.StatusCode, err
}

func (instance TestCSBInstance) BrokerUrl(subPath string) string {
	return fmt.Sprintf("http://localhost:%s/v2/%s", instance.port, subPath)
}

func BuildBrokerWithProvider(provider TerraformMock) TestCSBInstance {
	// build binary
	brokerPackDir := "/Users/jatinnaik/workspace/csb/csb-brokerpak-gcp/"
	csbBuild, err := gexec.Build("github.com/cloudfoundry-incubator/cloud-service-broker")
	Expect(err).NotTo(HaveOccurred())

	// create manifest.yml
	// we can replace the template, after this code is moved over the csb codebase,
	// when it has access to the internal manifest package
	workingDir, err := createWorkspace(brokerPackDir, provider.binary)
	Expect(err).NotTo(HaveOccurred())

	// build the broker pack
	command := exec.Command(csbBuild, "pak", "build")
	command.Dir = workingDir
	session, err := gexec.Start(command, GinkgoWriter, GinkgoWriter)

	Expect(err).NotTo(HaveOccurred())
	session.Wait(5 * time.Minute)

	return TestCSBInstance{brokerBuild: csbBuild, workspace: workingDir, username: "u", password: "p", port: "8080"}
}

func createWorkspace(brokerPackDir string, build string) (string, error) {
	workingDir, err := ioutil.TempDir("", "prefix")
	if err != nil {
		return "", err
	}
	err = linkBrokerpackFiles(brokerPackDir, workingDir)
	if err != nil {
		return "", err
	}

	return workingDir, templateManifest(brokerPackDir, build, workingDir)
}

func linkBrokerpackFiles(brokerPackDir string, workingDir string) error {
	yamlFiles, err := filepath.Glob(brokerPackDir + "*.yml")
	if err != nil {
		return err
	}
	for _, file := range yamlFiles {
		err = os.Link(file, path.Join(workingDir, filepath.Base(file)))
		if err != nil {
			return err
		}
	}
	err = os.Symlink(path.Join(brokerPackDir, "terraform"), path.Join(workingDir, "terraform"))
	if err != nil {
		return err
	}
	err = os.Remove(path.Join(workingDir, "manifest.yml"))
	if err != nil {
		return err
	}
	return nil
}

func templateManifest(brokerPackDir string, build string, workingDir string) error {
	contents, err := ioutil.ReadFile(path.Join(brokerPackDir, "manifest.yml"))
	if err != nil {
		return err
	}
	tmpl, err := template.New("config-template").Parse(string(contents))
	if err != nil {
		return err
	}
	outputFile, err := os.Create(path.Join(workingDir, "manifest.yml"))
	if err != nil {
		return err
	}

	err = tmpl.Execute(outputFile, struct{ Build string }{Build: build})

	if err != nil {
		return err
	}
	return outputFile.Close()
}
