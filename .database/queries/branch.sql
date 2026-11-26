-- name: CreateBranchV1 :one
SELECT *
FROM root.create_branch_v1(
    sqlc.arg('p_name'),
    sqlc.arg('p_address'),
    sqlc.narg('p_note')
);

-- name: GetBranchesV1 :many
SELECT *
FROM root.get_branches_v1(
    sqlc.narg('p_q'),
    sqlc.narg('p_limit'),
    sqlc.narg('p_offset')
);

-- name: CountBranchesV1 :one
SELECT root.count_branches_v1(sqlc.narg('p_q'));

-- name: GetBranchV1 :one
SELECT *
FROM root.get_branch_v1(sqlc.arg('p_id'));

-- name: UpdateBranchV1 :one
SELECT *
FROM root.update_branch_v1(
    sqlc.arg('p_id'),
    sqlc.arg('p_name'),
    sqlc.arg('p_address'),
    sqlc.narg('p_note'),
    sqlc.arg('p_activity_status')
);