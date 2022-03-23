provider "postgresql" {
  host      = var.hostname
  port      = local.port
  username  = var.admin_username
  password  = var.admin_password
  superuser = false
  database  = var.db_name
  sslmode   = var.use_tls ? "require" : "disable"
  clientcert {
      cert = "${path.module}/client_ca_cert.pem"
      key  = "${path.module}/client_private_key.pem"
  }
}

provider "local" {}