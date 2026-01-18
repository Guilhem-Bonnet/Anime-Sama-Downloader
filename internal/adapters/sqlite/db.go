package sqlite

import (
	"context"
	"database/sql"
	"embed"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"

	_ "modernc.org/sqlite"
)

//go:embed migrations/*.sql
var migrationsFS embed.FS

type DB struct {
	SQL *sql.DB
}

func Open(ctx context.Context, path string) (*DB, error) {
	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, err
	}

	// Base de sûreté, ajustable plus tard.
	db.SetMaxOpenConns(1)
	db.SetMaxIdleConns(1)
	db.SetConnMaxLifetime(0)

	ctxPing, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	if err := db.PingContext(ctxPing); err != nil {
		_ = db.Close()
		return nil, err
	}

	wrapper := &DB{SQL: db}
	if err := wrapper.Migrate(ctx); err != nil {
		_ = db.Close()
		return nil, err
	}
	return wrapper, nil
}

func (d *DB) Close() error {
	return d.SQL.Close()
}

func (d *DB) Migrate(ctx context.Context) error {
	if _, err := d.SQL.ExecContext(ctx, `CREATE TABLE IF NOT EXISTS schema_migrations (version INTEGER PRIMARY KEY, applied_at TEXT NOT NULL);`); err != nil {
		return err
	}

	applied, err := d.appliedVersions(ctx)
	if err != nil {
		return err
	}

	entries, err := migrationsFS.ReadDir("migrations")
	if err != nil {
		return err
	}

	versions := make([]int, 0, len(entries))
	byVersion := map[int]string{}
	for _, e := range entries {
		name := e.Name()
		if !strings.HasSuffix(name, ".sql") {
			continue
		}
		prefix := strings.SplitN(name, "_", 2)[0]
		v, err := strconv.Atoi(prefix)
		if err != nil {
			return fmt.Errorf("invalid migration name: %s", name)
		}
		versions = append(versions, v)
		byVersion[v] = name
	}
	sort.Ints(versions)

	for _, v := range versions {
		if applied[v] {
			continue
		}
		b, err := migrationsFS.ReadFile("migrations/" + byVersion[v])
		if err != nil {
			return err
		}
		upSQL := extractUp(string(b))
		if strings.TrimSpace(upSQL) == "" {
			continue
		}

		tx, err := d.SQL.BeginTx(ctx, nil)
		if err != nil {
			return err
		}
		if _, err := tx.ExecContext(ctx, upSQL); err != nil {
			_ = tx.Rollback()
			return fmt.Errorf("migration %s failed: %w", byVersion[v], err)
		}
		if _, err := tx.ExecContext(ctx, `INSERT INTO schema_migrations(version, applied_at) VALUES(?, ?)`, v, time.Now().UTC().Format(time.RFC3339)); err != nil {
			_ = tx.Rollback()
			return err
		}
		if err := tx.Commit(); err != nil {
			return err
		}
	}

	return nil
}

func (d *DB) appliedVersions(ctx context.Context) (map[int]bool, error) {
	rows, err := d.SQL.QueryContext(ctx, `SELECT version FROM schema_migrations`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out := map[int]bool{}
	for rows.Next() {
		var v int
		if err := rows.Scan(&v); err != nil {
			return nil, err
		}
		out[v] = true
	}
	return out, rows.Err()
}

func extractUp(sqlText string) string {
	lines := strings.Split(sqlText, "\n")
	var out []string
	inUp := false
	for _, line := range lines {
		trim := strings.TrimSpace(line)
		if strings.HasPrefix(trim, "-- +migrate Up") {
			inUp = true
			continue
		}
		if strings.HasPrefix(trim, "-- +migrate Down") {
			inUp = false
			continue
		}
		if inUp {
			out = append(out, line)
		}
	}
	return strings.Join(out, "\n")
}
