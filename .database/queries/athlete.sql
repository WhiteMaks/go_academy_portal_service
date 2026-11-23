-- name: CreateAthleteV1 :one
SELECT *
FROM root.create_athlete_v1(
	sqlc.arg('second_name'),
	sqlc.arg('first_name'),
	sqlc.narg('middle_name'),
	sqlc.arg('birth_date'),
	sqlc.narg('photo_url'),
	sqlc.narg('phone'),
	sqlc.narg('parent_second_name'),
	sqlc.narg('parent_first_name'),
	sqlc.narg('parent_middle_name'),
	sqlc.narg('parent_phone')
);

-- name: GetAthletesV1 :many
SELECT *
FROM root.get_athletes_v1();
