resource "google_sql_database_instance" "instance" {
  name             = var.instance_name
  database_version = var.database_version
  region           = var.region

  settings {
    tier        = local.service_tiers[var.cores]
    disk_size   = var.storage_gb
    user_labels = var.labels

    ip_configuration {
      ipv4_enabled    = false
      private_network = local.authorized_network_id
      #require_ssl = var.use_tls
    }
  }

  deletion_protection = false
}

resource "google_sql_database" "database" {
  name     = var.db_name
  instance = google_sql_database_instance.instance.name
}

resource "random_string" "username" {
  length  = 16
  special = false
}

resource "random_password" "password" {
  length           = 16
  special          = true
  override_special = "_@"
}

resource "google_sql_user" "admin_user" {
  name     = random_string.username.result
  instance = google_sql_database_instance.instance.name
  password = random_password.password.result
  deletion_policy="ABANDON"
}

# //create a non login user
resource "random_string" "nologin_username" {
  length  = 16
  special = false
}
resource "random_password" "nologin_password" {
  length           = 16
  special          = true
  override_special = "_@"
}

resource "postgresql_role" "non_login_user" {
  depends_on  = [google_sql_user.admin_user]
  name                = random_string.nologin_username.result
  password            = random_password.nologin_password.result
  login               = false
  skip_drop_role      = true
  skip_reassign_owned = true
}

# resource "null_resource" "create_nologin_role" {
#   depends_on  = [google_sql_user.admin_user]

#   provisioner "local-exec" {
#     command = format("psql -h %s -d %s -U %s --no-password -c \"CREATE ROLE nologin_role WITH PASSWORD='%s';\"",
#                      google_sql_database_instance.instance.first_ip_address,
#                      google_sql_database.database.name,
#                      google_sql_user.admin_user.name,
#                      random_string.nologin_password.result)
#     environment = {
#       PGPASSWORD = google_sql_user.admin_user.password
#     }
#   }

#   provisioner "local-exec" {
# 	when = destroy
#     command = format("psql -h %s -d %s -U %s --no-password -c \"DROP ROLE nologin_role;\"",
#                      google_sql_database_instance.instance.first_ip_address,
#                      google_sql_database.database.name,
#                      google_sql_user.admin_user.name)
#     environment = {
#       PGPASSWORD = google_sql_user.admin_user.password
#     }
#   }
# }