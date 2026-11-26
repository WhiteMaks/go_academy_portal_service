-- name: CreateCoachV1 :one
SELECT *
FROM root.create_coach_v1(
    sqlc.arg('p_last_name'),
    sqlc.arg('p_first_name'),
    sqlc.narg('p_patronymic'),
    sqlc.arg('p_phone'),
    sqlc.arg('p_specialization'),
    sqlc.narg('p_note'),
    sqlc.arg('p_birth_date')
);

-- name: GetCoachesV1 :many
SELECT *
FROM root.get_coaches_v1(
    sqlc.narg('p_q'),
    sqlc.narg('p_limit'),
    sqlc.narg('p_offset')
);

-- name: CountCoachesV1 :one
SELECT root.count_coaches_v1(sqlc.narg('p_q'));

-- name: GetCoachV1 :one
SELECT *
FROM root.get_coach_v1(sqlc.arg('p_id'));

-- name: UpdateCoachV1 :one
SELECT *
FROM root.update_coach_v1(
    sqlc.arg('p_id'),
    sqlc.arg('p_last_name'),
    sqlc.arg('p_first_name'),
    sqlc.narg('p_patronymic'),
    sqlc.arg('p_phone'),
    sqlc.arg('p_specialization'),
    sqlc.narg('p_note'),
    sqlc.arg('p_activity_status'),
    sqlc.arg('p_birth_date')
);