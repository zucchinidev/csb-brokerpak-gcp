provider "mysql" {
  endpoint = format("%s:%d", var.mysql_hostname, local.port)
  username = var.admin_username
  password = var.admin_password
  clientcert {
      cert = "${path.module}/client_ca_cert.pem"
      key  = "${path.module}/client_private_key.pem"
  }
}
