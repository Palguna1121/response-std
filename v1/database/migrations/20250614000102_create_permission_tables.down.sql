-- Drop permission tables in correct order to avoid foreign key constraints
DROP TABLE IF EXISTS role_has_permissions;
DROP TABLE IF EXISTS model_has_roles;
DROP TABLE IF EXISTS model_has_permissions;
DROP TABLE IF EXISTS roles;
DROP TABLE IF EXISTS permissions;