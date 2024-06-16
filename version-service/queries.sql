-- name: CreateNewVersion :one
INSERT INTO versions (service_id, image, tag, created_by)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: GetVersionsForService :many
SELECT *
FROM versions
WHERE service_id = $1
ORDER BY created_at DESC;

-- name: GetVersions :many
SELECT jsonb_object_agg(service_id, jsonb_build_object(
        'image', image,
        'tag', tag,
        'created_at', created_at,
        'created_by', created_by
    )) AS result
FROM (SELECT service_id,
             image,
             tag,
             created_at,
             created_by
      FROM versions v1
      WHERE created_at = (SELECT MAX(created_at)
                          FROM versions v2
                          WHERE v2.service_id = v1.service_id)) subquery;


