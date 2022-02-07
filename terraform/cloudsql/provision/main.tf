resource "google_sql_database_instance" "instance" {
  name             = var.instance_name
  database_version = var.database_version
  region           = var.region


  master_instance_name = var.master_instance_name
  #replica_configuration = var.replica_configuration
  root_password = var.root_password
  #restore_backup_context = var.restore_backup_context
  #clone = var.clone

  settings {
    tier        = local.service_tiers[var.cores]
    disk_size   = var.storage_gb
    user_labels = var.labels

    activation_policy = var.activation_policy
    availability_type = var.availability_type
    collation = var.collation
    disk_autoresize = var.disk_autoresize
    disk_type = var.disk_type
    pricing_plan = var.pricing_plan

    ip_configuration {
      ipv4_enabled    = false
      private_network = local.authorized_network_id
      #require_ssl = var.use_tls
    }

    backup_configuration {
        binary_log_enabled = var.binary_log_enabled
        enabled = var.enabled
        start_time = var.start_time
        point_in_time_recovery_enabled = var.point_in_time_recovery_enabled
        location = var.location
        transaction_log_retention_days = var.transaction_log_retention_days
    }

  }


  deletion_protection = var.deletion_protection
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
}
