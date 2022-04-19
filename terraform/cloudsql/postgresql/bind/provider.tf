provider "postgresql" {
  host      = var.hostname
  port      = local.port
  username  = var.admin_username
  password  = var.admin_password
  superuser = false
  database  = var.db_name
  sslmode   = var.use_tls ? "require" : "disable"
}

provider "csbpg" {
  host      = var.hostname
  port      = local.port
  username  = var.admin_username
  password  = var.admin_password
  database  = var.db_name
}
