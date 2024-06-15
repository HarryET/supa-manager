-- name: GetProjectsForAccountId :many
SELECT p.*
FROM organization_membership om
         JOIN project p on om.organization_id = p.organization_id
WHERE account_id = $1;