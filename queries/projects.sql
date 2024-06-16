-- name: GetProjectsForAccountId :many
SELECT p.*
FROM organization_membership om
         JOIN project p on om.organization_id = p.organization_id
WHERE account_id = $1;

-- name: CreateProject :one
INSERT INTO project (project_ref, project_name, organization_id, status, jwt_secret, cloud_provider, region)
VALUES ($1, $2, $3, 'UNKNOWN', $4, $5, $6)
RETURNING *;

-- name: GetProjectByRef :one
SELECT *
FROM project
WHERE project_ref = $1;