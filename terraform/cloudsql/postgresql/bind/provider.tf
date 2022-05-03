provider "csbpg" {
  host            = var.hostname
  port            = local.port
  username        = var.admin_username
  password        = var.admin_password
  database        = var.db_name
  data_owner_role = "data_owner_role"
  sslmode     = "verify-ca"
  sslrootcert = var.sslrootcert
  clientcert {
    cert = var.sslcert
    key  = var.sslkey
  }
}