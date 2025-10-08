-- name: IsAdminCreationAllowedV1 :one
SELECT * FROM root.is_admin_creation_allowed_v1();

-- name: CreateUserV1 :one
SELECT * FROM root.create_user_v1($1, $2, $3, $4);