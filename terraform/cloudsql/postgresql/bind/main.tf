resource "random_string" "username" {
  length  = 16
  special = false
  number  = false
}

resource "random_password" "password" {
  length           = 64
  override_special = "~_-."
  min_upper        = 2
  min_lower        = 2
  min_special      = 2
}

resource "postgresql_role" "new_user" {
  name                = random_string.username.result
  login               = true
  password            = random_password.password.result
  roles               = [
    var.admin_username
  ]
}

resource "local_file" "client_private_key" {
    content = var.client_private_key
    filename = "${path.module}/client_private_key.pem"
    file_permission = "0600"
}

resource "postgresql_role" "new_user" {
  name                = random_string.username.result
  login               = true
  password            = random_password.password.result
  skip_reassign_owned = true
  skip_drop_role      = true
}
