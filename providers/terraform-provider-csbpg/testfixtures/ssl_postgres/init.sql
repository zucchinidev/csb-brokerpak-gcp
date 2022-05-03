-- Creating role with same privilege as Google Super user role: https://cloud.google.com/sql/docs/postgres/users
CREATE ROLE cloudsqlsuperuser WITH LOGIN PASSWORD 'password' NOSUPERUSER CREATEDB CREATEROLE;
