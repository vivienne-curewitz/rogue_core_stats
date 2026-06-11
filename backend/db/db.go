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
			PRIMARY KEY (run_id, player_id, item_id, item_link)
		);`,
		`CREATE TABLE IF NOT EXISTS items (
			run_id TEXT NOT NULL,
			player_id TEXT NOT NULL,
			item_id TEXT NOT NULL,
			reference TEXT NOT NULL,
			PRIMARY KEY (run_id, player_id, item_id, reference)
		);`,
		`CREATE TABLE IF NOT EXISTS run_info (
			player_id TEXT NOT NULL,
			run_id TEXT NOT NULL,
			character_id TEXT NOT NULL,
			boss_id TEXT NOT NULL,
			status BOOLEAN NOT NULL,
			depth INTEGER NOT NULL DEFAULT 0,
			player_damage REAL NOT NULL DEFAULT 0,
			overkill_damage REAL NOT NULL DEFAULT 0,
			player_kills INTEGER NOT NULL DEFAULT 0,
			player_deaths INTEGER NOT NULL DEFAULT 0,
			total_stages INTEGER NOT NULL DEFAULT 0,
			completed_stages INTEGER NOT NULL DEFAULT 0,
			runtime INTEGER NOT NULL DEFAULT 0,
			player_rank INTEGER NOT NULL DEFAULT 0,
			character_rank INTEGER NOT NULL DEFAULT 0,
			character_stars INTEGER NOT NULL DEFAULT 0,
			minerals_mined REAL NOT NULL DEFAULT 0,
			max_armor REAL NOT NULL DEFAULT 0,
			max_health REAL NOT NULL DEFAULT 0,
			health_restored REAL NOT NULL DEFAULT 0,
			PRIMARY KEY (player_id, run_id)
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

func BatchWriteRunInfo(ctx context.Context, infos []types.RunOverview) error {
	batch := &pgx.Batch{}
	for _, info := range infos {
		batch.Queue(`INSERT INTO run_info (
			player_id, run_id, character_id, boss_id, status, depth, player_damage, overkill_damage, 
			player_kills, player_deaths, total_stages, completed_stages, runtime, player_rank, 
			character_rank, character_stars, minerals_mined, max_armor, max_health, health_restored
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20)
		ON CONFLICT (player_id, run_id) DO UPDATE SET 
			character_id = EXCLUDED.character_id, boss_id = EXCLUDED.boss_id, status = EXCLUDED.status,
			depth = EXCLUDED.depth, player_damage = EXCLUDED.player_damage, overkill_damage = EXCLUDED.overkill_damage,
			player_kills = EXCLUDED.player_kills, player_deaths = EXCLUDED.player_deaths, total_stages = EXCLUDED.total_stages,
			completed_stages = EXCLUDED.completed_stages, runtime = EXCLUDED.runtime, player_rank = EXCLUDED.player_rank,
			character_rank = EXCLUDED.character_rank, character_stars = EXCLUDED.character_stars,
			minerals_mined = EXCLUDED.minerals_mined, max_armor = EXCLUDED.max_armor,
			max_health = EXCLUDED.max_health, health_restored = EXCLUDED.health_restored`,
			info.PlayerId, info.RunId, info.CharacterId, info.BossId, info.Status, info.Depth, info.PlayerDamage, info.OverkillDamage,
			info.PlayerKills, info.PlayerDeaths, info.TotalStages, info.CompletedStages, info.Runtime, info.PlayerRank,
			info.CharacterRank, info.CharacterStars, info.MineralsMined, info.MaxArmor, info.MaxHealth, info.HealthRestored)
	}

	br := Pool.SendBatch(ctx, batch)
	defer br.Close()

	for range infos {
		_, err := br.Exec()
		if err != nil {
			return fmt.Errorf("failed to execute batch insert for run_info: %w", err)
		}
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
		`DROP TABLE IF EXISTS run_info;`,
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
                     ON CONFLICT (run_id, player_id, item_id, item_link) 
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
                     ON CONFLICT (run_id, player_id, item_id, reference) DO NOTHING`,
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

func GetUpgradesByRunID(ctx context.Context, runID string) ([]types.Upgrade, error) {
	rows, err := Pool.Query(ctx, `SELECT run_id, player_id, item_id, quantity, item_link FROM upgrades WHERE run_id = $1`, runID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var upgrades []types.Upgrade
	for rows.Next() {
		var u types.Upgrade
		if err := rows.Scan(&u.RunId, &u.PlayerId, &u.UpgradeId, &u.Quantity, &u.Reference); err != nil {
			return nil, err
		}
		upgrades = append(upgrades, u)
	}
	return upgrades, nil
}

func GetItemsByRunID(ctx context.Context, runID string) ([]types.Item, error) {
	rows, err := Pool.Query(ctx, `SELECT run_id, player_id, item_id, reference FROM items WHERE run_id = $1`, runID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []types.Item
	for rows.Next() {
		var i types.Item
		if err := rows.Scan(&i.RunId, &i.PlayerId, &i.ItemId, &i.Reference); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	return items, nil
}

func GetPlayerOverview(ctx context.Context, playerID string) ([]types.RunOverview, error) {
	query := `SELECT 
		player_id, run_id, character_id, boss_id, status, depth, player_damage, overkill_damage, 
		player_kills, player_deaths, total_stages, completed_stages, runtime, player_rank, 
		character_rank, character_stars, minerals_mined, max_armor, max_health, health_restored
	FROM run_info WHERE player_id = $1`
	rows, err := Pool.Query(ctx, query, playerID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch player overview: %w", err)
	}
	defer rows.Close()

	var overviews []types.RunOverview
	for rows.Next() {
		var o types.RunOverview
		if err := rows.Scan(
			&o.PlayerId, &o.RunId, &o.CharacterId, &o.BossId, &o.Status, &o.Depth, &o.PlayerDamage, &o.OverkillDamage,
			&o.PlayerKills, &o.PlayerDeaths, &o.TotalStages, &o.CompletedStages, &o.Runtime, &o.PlayerRank,
			&o.CharacterRank, &o.CharacterStars, &o.MineralsMined, &o.MaxArmor, &o.MaxHealth, &o.HealthRestored,
		); err != nil {
			return nil, fmt.Errorf("failed to scan run overview: %w", err)
		}
		overviews = append(overviews, o)
	}

	return overviews, nil
}
