-- name: GetMigration :one
SELECT * FROM public.migrations WHERE id = $1;

-- name: GetMigrations :many
SELECT * FROM public.migrations;

-- name: PutMigration :exec
INSERT INTO public.migrations (id, note) VALUES ($1, $2);