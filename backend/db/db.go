package db

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

var Pool *pgxpool.Pool

func InitDB(ctx context.Context) error {
	connStr := os.Getenv("DATABASE_URL")
	if connStr == "" {
		return fmt.Errorf("DATABASE_URL environment variable is not set")
	}

	var err error
	Pool, err = pgxpool.New(ctx, connStr)
	if err != nil {
		return fmt.Errorf("unable to create connection pool: %w", err)
	}

	if err := Pool.Ping(ctx); err != nil {
		return fmt.Errorf("unable to ping database: %w", err)
	}

	return createTables(ctx)
}

func createTables(ctx context.Context) error {
	queries := []string{
		`CREATE TABLE IF NOT EXISTS icon_to_uuid (
			uuid TEXT PRIMARY KEY,
			icon_name TEXT NOT NULL,
			icon_file TEXT NOT NULL
		);`,
		`CREATE TABLE IF NOT EXISTS upgrades (
			run_id TEXT NOT NULL,
			player_id TEXT NOT NULL,
			item_id TEXT NOT NULL,
			quantity INTEGER NOT NULL DEFAULT 1,
			item_link TEXT NOT NULL,
			PRIMARY KEY (run_id, player_id, item_id)
		);`,
		`CREATE TABLE IF NOT EXISTS items (
			run_id TEXT NOT NULL,
			player_id TEXT NOT NULL,
			item_id TEXT NOT NULL,
			reference TEXT NOT NULL,
			PRIMARY KEY (run_id, player_id, item_id)
		);`,
		`CREATE TABLE IF NOT EXISTS run_status (
			run_id TEXT NOT NULL,
			player_id TEXT NOT NULL,
			character_id TEXT NOT NULL,
			boss_id TEXT NOT NULL,
			status TEXT NOT NULL,
			PRIMARY KEY (run_id, player_id)
		);`,
		`CREATE TABLE IF NOT EXISTS run_damage (
			run_id TEXT NOT NULL,
			player_id TEXT NOT NULL,
			total_damage BIGINT NOT NULL DEFAULT 0,
			PRIMARY KEY (run_id, player_id)
		);`,
	}

	for _, q := range queries {
		if _, err := Pool.Exec(ctx, q); err != nil {
			return fmt.Errorf("failed to execute query %q: %w", q, err)
		}
	}

	return nil
}
