package conf

import (
	"context"
	"fmt"
	"github.com/harryet/supa-manager/database"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"os"
	"strings"
)

func EnsureMigrationsTableExists(conn *pgxpool.Pool) error {
	// Check if the table exists
	query := `
		SELECT EXISTS (
			SELECT FROM information_schema.tables 
			WHERE table_name = 'migrations'
		)
	`
	var exists bool
	err := conn.QueryRow(context.Background(), query).Scan(&exists)
	if err != nil {
		return fmt.Errorf("error checking if migrations table exists: %w", err)
	}

	// Create the table if it does not exist
	if !exists {
		createTableQuery := `
			CREATE TABLE IF NOT EXISTS public.migrations
			(
				id         text        not null primary key,
				note       text,
				applied_at timestamptz not null default now()
			)
		`
		_, err = conn.Exec(context.Background(), createTableQuery)
		if err != nil {
			return fmt.Errorf("error creating migrations table: %w", err)
		}
	}
	return nil
}

func EnsureMigrations(pool *pgxpool.Pool, conn *database.Queries) (bool, error) {
	migrations, err := conn.GetMigrations(context.Background())
	if err != nil {
		return false, err
	}

	files, err := os.ReadDir("migrations")
	if err != nil {
		return false, err
	}

	for _, file := range files {
		if file.IsDir() {
			continue
		}
		migrationId := strings.Split(file.Name(), "_")[0]
		migration, err := os.ReadFile(fmt.Sprintf("migrations/%s", file.Name()))
		if err != nil {
			return false, err
		}

		// check if the migration has been applied
		var applied bool
		for _, m := range migrations {
			if m.ID == migrationId {
				applied = true
				break
			}
		}

		if !applied {
			tx, err := pool.BeginTx(context.Background(), pgx.TxOptions{})
			if err != nil {
				return false, err
			}
			_, err = tx.Exec(context.Background(), string(migration))
			if err != nil {
				return false, err
			}
			defer tx.Rollback(context.Background())

			err = conn.WithTx(tx).PutMigration(context.Background(), database.PutMigrationParams{
				ID:   migrationId,
				Note: pgtype.Text{String: "applied automatically on start-up", Valid: true},
			})
			if err != nil {
				return false, err
			}
			if err = tx.Commit(context.Background()); err != nil {
				return false, err
			}
		}
	}

	return true, nil
}
