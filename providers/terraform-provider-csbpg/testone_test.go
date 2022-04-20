package main_test

import (
	"database/sql"
	"fmt"
	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	_ "github.com/lib/pq"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"
	"os/exec"
	"terraform-provider-csbpg/csbpg"
	"time"
)

const (
	username = "postgres"
	database = "postgres"
	hostname = "localhost"
)

var _ = Describe("Tests", func() {
	It("creates a binding user", func() {
		// Start the postgres docker image
		password := uuid.New().String()
		port := freePort()
		cmd := exec.Command(
			"docker", "run",
			"-e", fmt.Sprintf("POSTGRES_PASSWORD=%s", password),
			"-p", fmt.Sprintf("%d:5432", port),
			"-t", "postgres",
		)
		session, err := gexec.Start(cmd, GinkgoWriter, GinkgoWriter)
		Expect(err).NotTo(HaveOccurred())

		// Wait for the port open
		uri := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable", username, password, hostname, port, database)
		ping := func() error {
			db, err := sql.Open("postgres", uri)
			if err != nil {
				return err
			}
			defer db.Close()
			return db.Ping()
		}
		Eventually(ping).WithTimeout(30 * time.Second).WithPolling(time.Second).Should(Succeed())

		dataOwnerRole := uuid.New().String()
		bindingUsername := uuid.New().String()
		bindingPassword := uuid.New().String()
		hcl := fmt.Sprintf(`
		provider "csbpg" {
		  host            = "%s"
		  port            = %d
		  username        = "%s"
		  password        = "%s"
		  database        = "%s"
		  data_owner_role = "%s"
		}

		resource "csbpg_binding_user" "binding_user" {
		  username = "%s"
		  password = "%s"
		}
		`, hostname, port, username, password, database, dataOwnerRole, bindingUsername, bindingPassword)

		//	1. configure provider to connect to our docker
		// 2. Terraform apply
		resource.Test(GinkgoT(), resource.TestCase{
			IsUnitTest: true, // means we don't need to set TF_ACC
			ProviderFactories: map[string]func() (*schema.Provider, error){
				"csbpg": func() (*schema.Provider, error) { return csbpg.Provider(), nil },
			},
			Steps: []resource.TestStep{{
				ResourceName: "csbpg_shared_role.shared_role",
				Config:       hcl,
			}},
		})

		db, err := sql.Open("postgres", uri)
		Expect(err).NotTo(HaveOccurred())

		By("checking that the data owner role is created")
		rows, err := db.Query(fmt.Sprintf("SELECT FROM pg_catalog.pg_roles WHERE rolname = '%s'", dataOwnerRole))
		Expect(err).NotTo(HaveOccurred())
		Expect(rows.Next()).To(BeTrue(), fmt.Sprintf("role %q has not been created", dataOwnerRole))

		By("checking that the binding user is created")
		rows, err = db.Query(fmt.Sprintf("SELECT FROM pg_catalog.pg_roles WHERE rolname = '%s'", bindingUsername))
		Expect(err).NotTo(HaveOccurred())
		Expect(rows.Next()).To(BeTrue(), fmt.Sprintf("role %q has not been created", bindingUsername))

		By("checking that the binding user is a member of the data owner role")
		rows, err = db.Query(fmt.Sprintf("SELECT pg_has_role('%s', '%s', 'member')", bindingUsername, dataOwnerRole))
		Expect(err).NotTo(HaveOccurred())
		var result bool
		Expect(rows.Next()).To(BeTrue(), "pg_has_role() query failed")
		Expect(rows.Scan(&result)).To(Succeed())
		Expect(result).To(BeTrue(), "binding user is not a member of the data_owner_role")

		// check that binding creds work

		// 1. Connect to the docker instance
		// 2. bind port to local

		// get the host
		// 3. Assert the users is created by connecting the docker host using postgres

		// Kill the docker container
		session.Terminate()
	})
})
