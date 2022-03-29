output "name" { value = google_sql_database.database.name }
output "hostname" { value = google_sql_database_instance.instance.first_ip_address }

output "username" { value = postgresql_role.createrole_user.name }
output "password" {
  sensitive = true
  value     = postgresql_role.createrole_user.password
}
output "use_tls" { value = var.use_tls }

output "client_ca_cert" { value = google_sql_ssl_cert.client_cert.cert }
output "client_private_key" {
    value = google_sql_ssl_cert.client_cert.private_key
    sensitive   = true
}