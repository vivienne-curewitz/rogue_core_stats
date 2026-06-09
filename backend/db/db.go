package db

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/vivienne-curewitz/rogue_core_stats/types"
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
		`CREATE TABLE IF NOT EXISTS run_info (
			run_id TEXT NOT NULL,
			player_id TEXT NOT NULL,
			character_id TEXT NOT NULL,
			boss_id TEXT NOT NULL,
			status TEXT NOT NULL,
			PRIMARY KEY (run_id, player_id)
		);`,
		`CREATE TABLE IF NOT EXISTS run_status (
			run_id TEXT PRIMARY KEY,
			status BOOLEAN NOT NULL
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

func WriteRunStatus(ctx context.Context, status types.RunStatus) error {
	query := `INSERT INTO run_status (run_id, status) VALUES ($1, $2)
              ON CONFLICT (run_id) DO UPDATE SET status = EXCLUDED.status`
	_, err := Pool.Exec(ctx, query, status.RunId, status.Status)
	if err != nil {
		return fmt.Errorf("failed to write run status: %w", err)
	}
	return nil
}

func RunExists(ctx context.Context, runId string) (bool, error) {
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM run_status WHERE run_id = $1)`
	err := Pool.QueryRow(ctx, query, runId).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to check if run exists: %w", err)
	}
	return exists, nil
}

// DebugDropTables drops the run_status, items, and upgrades tables.
// This should only be used for testing and debugging.
func DebugDropTables(ctx context.Context) error {
	queries := []string{
		`DROP TABLE IF EXISTS run_status;`,
		`DROP TABLE IF EXISTS items;`,
		`DROP TABLE IF EXISTS upgrades;`,
	}

	for _, q := range queries {
		if _, err := Pool.Exec(ctx, q); err != nil {
			return fmt.Errorf("failed to drop table: %w", err)
		}
	}
	return nil
}

func BatchWriteUpgrades(ctx context.Context, upgrades []types.Upgrade) error {
	batch := &pgx.Batch{}
	for _, u := range upgrades {
		batch.Queue(`INSERT INTO upgrades (run_id, player_id, item_id, quantity, item_link) 
                     VALUES ($1, $2, $3, $4, $5) 
                     ON CONFLICT (run_id, player_id, item_id) 
                     DO UPDATE SET quantity = upgrades.quantity + EXCLUDED.quantity`,
			u.RunId, u.PlayerId, u.UpgradeId, u.Quantity, u.Reference)
	}

	br := Pool.SendBatch(ctx, batch)
	defer br.Close()

	for range upgrades {
		_, err := br.Exec()
		if err != nil {
			return fmt.Errorf("failed to execute batch insert for upgrades: %w", err)
		}
	}
	return nil
}

func BatchWriteItems(ctx context.Context, items []types.Item) error {
	batch := &pgx.Batch{}
	for _, i := range items {
		batch.Queue(`INSERT INTO items (run_id, player_id, item_id, reference) 
                     VALUES ($1, $2, $3, $4) 
                     ON CONFLICT (run_id, player_id, item_id) DO NOTHING`,
			i.RunId, i.PlayerId, i.ItemId, i.Reference)
	}

	br := Pool.SendBatch(ctx, batch)
	defer br.Close()

	for range items {
		_, err := br.Exec()
		if err != nil {
			return fmt.Errorf("failed to execute batch insert for items: %w", err)
		}
	}
	return nil
}
