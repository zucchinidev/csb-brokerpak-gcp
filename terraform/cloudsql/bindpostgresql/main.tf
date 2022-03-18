resource "random_string" "username" {
  length  = 16
  special = false
  number  = false
}

resource "random_password" "user_password" {
  length           = 64
  override_special = "~_-."
  min_upper        = 2
  min_lower        = 2
  min_special      = 2
}

resource "random_password" "nologin_password" {
  length           = 64
  override_special = "~_-."
  min_upper        = 2
  min_lower        = 2
  min_special      = 2
}

resource "postgresql_role" "new_user" {
  name                = random_string.username.result
  password            = random_password.user_password.result
  login               = true
  roles               = [
    var.nologin_user_role
  ]
}

resource "postgresql_grant" "db_access" {
  depends_on  = [postgresql_role.new_user]
  database    = var.db_name
  role        = postgresql_role.new_user.name
  object_type = "database"
  privileges  = ["ALL"]
}

resource "postgresql_grant" "table_access" {
  depends_on  = [postgresql_role.new_user]
  database    = var.db_name
  role        = postgresql_role.new_user.name
  schema      = "public"
  object_type = "table"
  privileges  = ["ALL"]
}
