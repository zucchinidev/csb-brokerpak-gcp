output "username" { value = random_string.username.result }
output "password" {
  sensitive = true
  value     = random_password.user_password.result
}
output "uri" {
  sensitive = true
  value = format("postgresql://%s:%s@%s:%d/%s",
    random_string.username.result,
    random_password.user_password.result,
    var.hostname,
    var.port,
  var.db_name)
}

output "jdbcUrl" {
  sensitive = true
  value = format("jdbc:%s://%s:%s/%s?user=%s\u0026password=%s\u0026verifyServerCertificate=true\u0026useSSL=%v\u0026requireSSL=false",
    "postgresql",
    var.hostname,
    var.port,
    var.db_name,
    random_string.username.result,
    random_password.user_password.result,
  var.use_tls)
}

# output "nologin_id" {
#   value     = "nologin_role"
# }