// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0
// source: migrations.sql

package database

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

const getMigration = `-- name: GetMigration :one
SELECT id, note, applied_at FROM public.migrations WHERE id = $1
`

func (q *Queries) GetMigration(ctx context.Context, id string) (Migration, error) {
	row := q.db.QueryRow(ctx, getMigration, id)
	var i Migration
	err := row.Scan(&i.ID, &i.Note, &i.AppliedAt)
	return i, err
}

const getMigrations = `-- name: GetMigrations :many
SELECT id, note, applied_at FROM public.migrations
`

func (q *Queries) GetMigrations(ctx context.Context) ([]Migration, error) {
	rows, err := q.db.Query(ctx, getMigrations)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Migration
	for rows.Next() {
		var i Migration
		if err := rows.Scan(&i.ID, &i.Note, &i.AppliedAt); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const putMigration = `-- name: PutMigration :exec
INSERT INTO public.migrations (id, note) VALUES ($1, $2)
`

type PutMigrationParams struct {
	ID   string
	Note pgtype.Text
}

func (q *Queries) PutMigration(ctx context.Context, arg PutMigrationParams) error {
	_, err := q.db.Exec(ctx, putMigration, arg.ID, arg.Note)
	return err
}
