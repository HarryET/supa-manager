// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0
// source: projects.sql

package database

import (
	"context"
)

const createProject = `-- name: CreateProject :one
INSERT INTO project (project_ref, project_name, organization_id, status, jwt_secret)
VALUES ($1, $2, $3, 'UNKNOWN', $4)
RETURNING id, project_ref, project_name, organization_id, status, cloud_provider, region, jwt_secret, created_at, updated_at
`

type CreateProjectParams struct {
	ProjectRef     string
	ProjectName    string
	OrganizationID int32
	JwtSecret      string
}

func (q *Queries) CreateProject(ctx context.Context, arg CreateProjectParams) (Project, error) {
	row := q.db.QueryRow(ctx, createProject,
		arg.ProjectRef,
		arg.ProjectName,
		arg.OrganizationID,
		arg.JwtSecret,
	)
	var i Project
	err := row.Scan(
		&i.ID,
		&i.ProjectRef,
		&i.ProjectName,
		&i.OrganizationID,
		&i.Status,
		&i.CloudProvider,
		&i.Region,
		&i.JwtSecret,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const getProjectByRef = `-- name: GetProjectByRef :one
SELECT id, project_ref, project_name, organization_id, status, cloud_provider, region, jwt_secret, created_at, updated_at
FROM project
WHERE project_ref = $1
`

func (q *Queries) GetProjectByRef(ctx context.Context, projectRef string) (Project, error) {
	row := q.db.QueryRow(ctx, getProjectByRef, projectRef)
	var i Project
	err := row.Scan(
		&i.ID,
		&i.ProjectRef,
		&i.ProjectName,
		&i.OrganizationID,
		&i.Status,
		&i.CloudProvider,
		&i.Region,
		&i.JwtSecret,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const getProjectsForAccountId = `-- name: GetProjectsForAccountId :many
SELECT p.id, p.project_ref, p.project_name, p.organization_id, p.status, p.cloud_provider, p.region, p.jwt_secret, p.created_at, p.updated_at
FROM organization_membership om
         JOIN project p on om.organization_id = p.organization_id
WHERE account_id = $1
`

func (q *Queries) GetProjectsForAccountId(ctx context.Context, accountID int32) ([]Project, error) {
	rows, err := q.db.Query(ctx, getProjectsForAccountId, accountID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Project
	for rows.Next() {
		var i Project
		if err := rows.Scan(
			&i.ID,
			&i.ProjectRef,
			&i.ProjectName,
			&i.OrganizationID,
			&i.Status,
			&i.CloudProvider,
			&i.Region,
			&i.JwtSecret,
			&i.CreatedAt,
			&i.UpdatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
