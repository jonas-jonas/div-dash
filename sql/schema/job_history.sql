-- name: GetJob :one
SELECT *
FROM "job_history"
WHERE id = $1;

-- name: GetJobsByName :many
SELECT *
FROM "job_history"
WHERE name = $1
ORDER BY started;

-- name: GetLastJobByName :one
SELECT *, CASE WHEN error_message IS NULL THEN false ELSE true END AS had_error
FROM "job_history"
WHERE name = $1
ORDER BY started DESC
LIMIT 1;

-- name: StartJob :one
INSERT INTO "job_history" (name, started)
VALUES ($1, $2)
RETURNING id, started;

-- name: FinishJob :one
UPDATE "job_history"
SET "finished" = $1, "error_message" = $2
WHERE id = $3
RETURNING "name", "id", "started", "finished";
