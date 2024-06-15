-- name: CreateOrganizationMembership :one
INSERT INTO organization_membership (organization_id, account_id, role)
VALUES ($1, $2, $3)
RETURNING *;