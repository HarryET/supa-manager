package utils

import (
	"github.com/jackc/pgx/v5/pgtype"
)

func PgTextToPointer(ns pgtype.Text) *string {
	if ns.Valid {
		return &ns.String
	}

	return nil
}
