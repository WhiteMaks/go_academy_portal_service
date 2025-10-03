-- name: IsAdminCreationAllowed :one
SELECT * FROM root.is_admin_creation_allowed();

-- name: CreateUser :one
SELECT * FROM root.create_user($1, $2, $3, $4, $5);