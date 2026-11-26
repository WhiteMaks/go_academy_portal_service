-- name: CreateHallV1 :one
SELECT *
FROM root.create_hall_v1(
    sqlc.arg('p_branch_id'),
    sqlc.arg('p_name'),
    sqlc.narg('p_subtitle'),
    sqlc.arg('p_capacity'),
    sqlc.narg('p_note')
);

-- name: GetBranchHallsV1 :many
SELECT *
FROM root.get_branch_halls_v1(
    sqlc.arg('p_branch_id'),
    sqlc.narg('p_q'),
    sqlc.narg('p_limit'),
    sqlc.narg('p_offset')
);

-- name: CountBranchHallsV1 :one
SELECT root.count_branch_halls_v1(
    sqlc.arg('p_branch_id'),
    sqlc.narg('p_q')
);

-- name: GetHallV1 :one
SELECT *
FROM root.get_hall_v1(sqlc.arg('p_id'));

-- name: UpdateHallV1 :one
SELECT *
FROM root.update_hall_v1(
    sqlc.arg('p_id'),
    sqlc.arg('p_name'),
    sqlc.narg('p_subtitle'),
    sqlc.arg('p_capacity'),
    sqlc.narg('p_note'),
    sqlc.arg('p_activity_status')
);