terraform {
  required_providers {
    csbpg = {
      source  = "cloudfoundry.org/cloud-service-broker/csbpg"
      version = "1.0.0"
    }
  }
}

provider "csbpg" {
  host = "localhost"
  port = 5432
  username = "postgres"
  password = "flubber"
  database = "postgres"
}

resource "csbpg_binding_user" "binding_user" {
  username = "bumble"
  password = "womble"
  shared_role = "flibble"
}
