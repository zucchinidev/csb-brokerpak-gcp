provider "postgresql" {
  host      = var.hostname
  port      = local.port
  username  = var.admin_username
  password  = var.admin_password
  superuser = false
  database  = var.db_name
  sslmode   = var.use_tls ? "require" : "disable"
  clientcert {
      cert = "${path.module}/sslcert.pem"
      key  = "${path.module}/sslkey.pem"
  }
}

provider "local" {}