package saveparser

import (
	"context"
	"fmt"
	"os"
	"sort"
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

	data, _ := os.ReadFile("example.sav")

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
	data, _ := os.ReadFile("example.sav")
	jsonStr, err := ConvertUesaveToJSON(data)
	if err != nil {
		t.Fatalf("Failed to convert save to JSON: %v", err)
	}

	runs := GetRunHistoryEntries(jsonStr)
	if len(runs) == 0 {
		t.Fatal("No runs found in history")
	}

	runString := runs[0]
	runId := GetRunID(runString)
	players := GetRunPlayers(runString)

	// Local extraction for "expected" data
	var expectedUpgrades []types.Upgrade
	var expectedItems []types.Item

	for _, player := range players {
		pUpgrades := GetRunUpgrades(player, runId)
		pItems, pItemUpgrades := GetRunItems(player, runId)

		allPlayerUpgrades := append(pUpgrades, pItemUpgrades...)
		allPlayerUpgrades = consolidateUpgrades(allPlayerUpgrades)

		expectedUpgrades = append(expectedUpgrades, allPlayerUpgrades...)
		expectedItems = append(expectedItems, pItems...)
	}

	// Write to DB using production code
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

	// Structural verification helper
	getFingerprints := func(items []types.Item, upgrades []types.Upgrade) map[string][]string {
		// Map Player -> list of item fingerprints
		playerItems := make(map[string][]string)

		// Group upgrades by (Player, Reference)
		type upgKey struct {
			PlayerID  string
			Reference string
		}
		upgGroups := make(map[upgKey][]string)
		playerLevelUpgrades := make(map[string][]string)

		for _, u := range upgrades {
			if u.Reference == "" {
				playerLevelUpgrades[u.PlayerId] = append(playerLevelUpgrades[u.PlayerId], fmt.Sprintf("%s:%d", u.UpgradeId, u.Quantity))
				continue
			}
			key := upgKey{u.PlayerId, u.Reference}
			upgGroups[key] = append(upgGroups[key], fmt.Sprintf("%s:%d", u.UpgradeId, u.Quantity))
		}

		// Sort player level upgrades
		for p := range playerLevelUpgrades {
			sort.Strings(playerLevelUpgrades[p])
		}

		// Create item fingerprints
		for _, it := range items {
			key := upgKey{it.PlayerId, it.Reference}
			uList := upgGroups[key]
			sort.Strings(uList)
			fingerprint := fmt.Sprintf("%s[%v]", it.ItemId, uList)
			playerItems[it.PlayerId] = append(playerItems[it.PlayerId], fingerprint)
		}

		// Sort item lists per player
		for p := range playerItems {
			sort.Strings(playerItems[p])
		}

		// Add player-level upgrades as a special "item"
		for p, ups := range playerLevelUpgrades {
			playerItems[p] = append(playerItems[p], fmt.Sprintf("PLAYER_LEVEL[%v]", ups))
			sort.Strings(playerItems[p])
		}

		return playerItems
	}

	expectedFingerprints := getFingerprints(expectedItems, expectedUpgrades)
	dbFingerprints := getFingerprints(dbItems, dbUpgrades)

	// Compare
	if len(expectedFingerprints) != len(dbFingerprints) {
		t.Errorf("Player count mismatch: expected %d, got %d", len(expectedFingerprints), len(dbFingerprints))
	}

	for p, eList := range expectedFingerprints {
		dList, ok := dbFingerprints[p]
		if !ok {
			t.Errorf("Player %s missing from DB", p)
			continue
		}
		if len(eList) != len(dList) {
			t.Errorf("Fingerprint count mismatch for player %s: expected %d, got %d", p, len(eList), len(dList))
			t.Logf("Expected: %v", eList)
			t.Logf("Got:      %v", dList)
			continue
		}
		for i := range eList {
			if eList[i] != dList[i] {
				t.Errorf("Fingerprint mismatch for player %s at index %d: expected %s, got %s", p, i, eList[i], dList[i])
			}
		}
	}
}
