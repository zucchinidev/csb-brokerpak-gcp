output "name" { value = google_sql_database.database.name }
output "hostname" { value = google_sql_database_instance.instance.first_ip_address }

output "port" { value = (var.database_version == "POSTGRES_11" ? 5432 : 3306) }
output "username" { value = postgresql_role.createrole_user.name }
output "password" {
  sensitive = true
  value     = postgresql_role.createrole_user.password
}
output "use_tls" { value = false }

output "nologin_user_role" { value = postgresql_role.createrole_user.id }
