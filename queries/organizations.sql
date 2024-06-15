-- name: GetOrganizationById :one
SELECT * FROM public.organizations WHERE slug = sqlc.arg('id');

