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
  name     = random_string.username.result
  login    = true
  password = random_password.password.result
  roles = [
    var.admin_username
  ]
}

resource "csbpg_shared_role" "shared_role" {
  name = "fakesharedrole"
}

resource "csbpg_binding_user" "binding_user" {
  username = "fakeusername"
  password = "fakepassword"
  shared_role = csbpg_shared_role.shared_role.name
}
