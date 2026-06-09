package saveparser

import (
	"context"
	"testing"

	"github.com/vivienne-curewitz/rogue_core_stats/db"
	"github.com/vivienne-curewitz/rogue_core_stats/types"
)

func TestExtractAndWrite(t *testing.T) {
	ctx := context.Background()

	// Initialize database for the test
	if err := db.LoadConfig(); err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}
	if err := db.InitDB(ctx); err != nil {
		t.Fatalf("Failed to init DB: %v", err)
	}

	// Drop tables for a clean test run
	if err := db.DebugDropTables(ctx); err != nil {
		t.Fatalf("Failed to drop tables: %v", err)
	}

	// Re-initialize to recreate the dropped tables
	if err := db.InitDB(ctx); err != nil {
		t.Fatalf("Failed to recreate tables: %v", err)
	}

	data, _ := readFile("example.sav")

	jsonStr, err := ConvertUesaveToJSON(data)
	if err != nil {
		t.Fatalf("Failed to convert save to JSON: %v", err)
	}

	runs := GetRunHistoryEntries(jsonStr)
	for _, run := range runs {
		inserted, err := ExtractRunData(run)
		if !inserted {
			t.Fatalf("Failed to insert data")
		}
		if err != nil {
			t.Fatalf("Failed with error: %v\n", err)
		}
	}

	inserted, err := ExtractRunData(runs[0])
	if inserted {
		t.Fatalf("Wrote duplicated run")
	}
	if err != nil {
		t.Fatalf("Failed with error: %v\n", err)
	}
}

func TestDatabaseIntegrity(t *testing.T) {
	ctx := context.Background()
	// Setup DB
	if err := db.LoadConfig(); err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}
	if err := db.InitDB(ctx); err != nil {
		t.Fatalf("Failed to init DB: %v", err)
	}
	if err := db.DebugDropTables(ctx); err != nil {
		t.Fatalf("Failed to drop tables: %v", err)
	}
	if err := db.InitDB(ctx); err != nil {
		t.Fatalf("Failed to recreate tables: %v", err)
	}

	// Read file
	data, _ := readFile("example.sav")
	jsonStr, err := ConvertUesaveToJSON(data)
	if err != nil {
		t.Fatalf("Failed to convert save to JSON: %v", err)
	}

	runs := GetRunHistoryEntries(jsonStr)
	if len(runs) == 0 {
		t.Fatal("No runs found in history")
	}

	// We'll perform extraction ONLY ONCE to ensure consistent UUIDs
	runString := runs[0]
	runId := GetRunID(runString)
	players := GetRunPlayers(runString)

	var expectedUpgrades []types.Upgrade
	var expectedItems []types.Item

	for _, player := range players {
		pUpgrades := GetRunUpgrades(player, runId)
		pItems, pItemUpgrades := GetRunItems(player, runId)

		allPlayerUpgrades := append(pUpgrades, pItemUpgrades...)
		// Consolidate per player as ExtractRunData does
		allPlayerUpgrades = consolidateUpgrades(allPlayerUpgrades)

		expectedUpgrades = append(expectedUpgrades, allPlayerUpgrades...)
		expectedItems = append(expectedItems, pItems...)
	}

	ExtractRunData(runString)

	// Fetch from DB
	dbUpgrades, err := db.GetUpgradesByRunID(ctx, runId)
	if err != nil {
		t.Fatalf("Failed to fetch upgrades from DB: %v", err)
	}
	dbItems, err := db.GetItemsByRunID(ctx, runId)
	if err != nil {
		t.Fatalf("Failed to fetch items from DB: %v", err)
	}

	// Compare lengths
	if len(expectedUpgrades) != len(dbUpgrades) {
		t.Errorf("Upgrades length mismatch: expected %d, got %d", len(expectedUpgrades), len(dbUpgrades))
	}
	if len(expectedItems) != len(dbItems) {
		t.Errorf("Items length mismatch: expected %d, got %d", len(expectedItems), len(dbItems))
	}

	// Helper to create maps for easy comparison
	checkUpgrades := make(map[string]types.Upgrade)
	for _, u := range dbUpgrades {
		// Key must be unique: player + upgrade + reference
		key := u.PlayerId + "|" + u.UpgradeId + "|" + u.Reference
		checkUpgrades[key] = u
	}

	for _, u := range expectedUpgrades {
		key := u.PlayerId + "|" + u.UpgradeId + "|" + u.Reference
		dbU, ok := checkUpgrades[key]
		if !ok {
			t.Errorf("Upgrade for player %s, id %s, ref %s missing from DB", u.PlayerId, u.UpgradeId, u.Reference)
			continue
		}
		if dbU.Quantity != u.Quantity {
			t.Errorf("Quantity mismatch for upgrade %s: expected %d, got %d", u.UpgradeId, u.Quantity, dbU.Quantity)
		}
	}

	checkItems := make(map[string]types.Item)
	for _, i := range dbItems {
		key := i.PlayerId + "|" + i.ItemId + "|" + i.Reference
		checkItems[key] = i
	}

	for _, i := range expectedItems {
		key := i.PlayerId + "|" + i.ItemId + "|" + i.Reference
		if _, ok := checkItems[key]; !ok {
			t.Errorf("Item for player %s, id %s, ref %s missing from DB", i.PlayerId, i.ItemId, i.Reference)
		}
	}
}
