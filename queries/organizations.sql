-- name: GetOrganizationById :one
SELECT * FROM public.organizations WHERE slug = sqlc.arg('id');

-- name: GetOrganizationsForAccountId :many
SELECT o.*, om.role as member_role
FROM organization_membership om
         JOIN organizations o on o.id = om.organization_id
WHERE account_id = $1;