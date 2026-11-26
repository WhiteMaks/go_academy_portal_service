-- name: CreateAthleteV1 :one
SELECT *
FROM root.create_athlete_v1(
    sqlc.arg('p_last_name'),
    sqlc.arg('p_first_name'),
    sqlc.narg('p_patronymic'),
    sqlc.arg('p_birth_date'),
    sqlc.arg('p_phone'),
    sqlc.narg('p_parent_last_name'),
    sqlc.narg('p_parent_first_name'),
    sqlc.narg('p_parent_patronymic'),
    sqlc.narg('p_parent_birth_date'),
    sqlc.narg('p_parent_phone'),
    sqlc.narg('p_note')
);

-- name: GetAthletesV1 :many
SELECT *
FROM root.get_athletes_v1(
    sqlc.narg('p_q'),
    sqlc.narg('p_limit'),
    sqlc.narg('p_offset')
);

-- name: CountAthletesV1 :one
SELECT root.count_athletes_v1(sqlc.narg('p_q'));

-- name: GetAthleteV1 :one
SELECT *
FROM root.get_athlete_v1(sqlc.arg('p_id'));