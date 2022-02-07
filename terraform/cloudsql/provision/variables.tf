variable "cores" { type = number }
variable "authorized_network" { type = string }
variable "authorized_network_id" { type = string }
variable "instance_name" { type = string }
variable "db_name" { type = string }
variable "deletion_protection" { type = bool }
variable "activation_policy" { type = string }
variable "availability_type" { type = string }
variable "collation" { type = string }
variable "disk_autoresize" { type = bool }
variable "disk_size" { type = number }
variable "disk_type" { type = string }
variable "pricing_plan" { type = string }

variable "master_instance_name" { type = string }
#variable "replica_configuration" { type = map(any) }
variable "root_password" { type = string }
#variable "restore_backup_context" { type = map(any) }
#variable "clone" { type = map(any) }


variable "binary_log_enabled" { type = bool }
variable "enabled" { type = bool }
variable "start_time" { type = string }
variable "point_in_time_recovery_enabled" { type = bool }
variable "location" { type = string }
variable "transaction_log_retention_days" { type = number }

variable "region" { type = string }
variable "labels" { type = map(any) }
variable "storage_gb" { type = number }
variable "database_version" { type = string }

variable "credentials" { type = string }
variable "project" { type = string }
#variable use_tls { type = bool }