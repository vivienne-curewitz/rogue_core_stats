package saveparser

import (
	"encoding/json"
	"log"
	"os"
	"testing"
)

func readFile(filePath string) ([]byte, error) {
	fileBytes, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	return fileBytes, nil
}

func TestSaveParser(t *testing.T) {
	filePath := "example.sav"
	data, err := readFile(filePath)
	if err != nil {
		t.Errorf("Could not load testing file: %s\n", err)
	}
	_, err = ConvertUesaveToJSON(data)
	if err != nil {
		t.Errorf("Failed to parse file: %s\n", err)
	}
}

func TestSaveParserStructure(t *testing.T) {
	filePath := "example.sav"
	fdata, err := readFile(filePath)
	if err != nil {
		t.Errorf("Could not load testing file: %s\n", err)
	}
	jsonStr, err := ConvertUesaveToJSON(fdata)
	if err != nil {
		t.Fatalf("Failed to parse file: %s\n", err)
	}

	var data map[string]interface{}
	err = json.Unmarshal([]byte(jsonStr), &data)
	if err != nil {
		t.Fatalf("Failed to unmarshal JSON: %s\n", err)
	}

	expectedKeys := []string{"header", "schemas", "root", "extra"}
	for _, key := range expectedKeys {
		if _, ok := data[key]; !ok {
			keys := make([]string, 0, len(data))
			for k := range data {
				keys = append(keys, k)
			}
			t.Errorf("Expected key %q not found in JSON top level. Actual keys: %v", key, keys)
		}
	}
}

func TestGetRunHistoryData(t *testing.T) {
	filePath := "example.sav"
	fdata, err := readFile(filePath)
	if err != nil {
		t.Errorf("Could not load testing file: %s\n", err)
	}
	jsonStr, err := ConvertUesaveToJSON(fdata)
	if err != nil {
		t.Fatalf("Failed to parse file: %s\n", err)
	}

	runHistoryData := GetRunHistoryEntries(jsonStr)
	if runHistoryData == nil {
		t.Errorf("Expected run history data, but got an empty string")
	}
	if len(runHistoryData) != 8 {
		t.Errorf("Expected 8 Runs, Got %d\n", len(runHistoryData))
	}
}

func TestGetRunHistoryIDs(t *testing.T) {
	filePath := "example.sav"
	fdata, err := readFile(filePath)
	if err != nil {
		t.Errorf("Could not load testing file: %s\n", err)
	}
	jsonStr, err := ConvertUesaveToJSON(fdata)
	if err != nil {
		t.Fatalf("Failed to parse file: %s\n", err)
	}

	runHistoryData := GetRunHistoryEntries(jsonStr)
	if runHistoryData == nil {
		t.Errorf("Expected run history data, but got an empty string")
	}
	if len(runHistoryData) != 8 {
		t.Errorf("Expected 8 Runs, Got %d\n", len(runHistoryData))
	}
	id_0 := GetRunID(runHistoryData[0])
	id_1 := GetRunID(runHistoryData[1])
	if id_0 == id_1 {
		t.Errorf("Expected different IDs for different runs, but got the same ID: %s", id_0)
	}
	id_0_repeated := GetRunID(runHistoryData[0])
	if id_0 != id_0_repeated {
		t.Errorf("Expected same ID for the same run, but got different IDs: %s and %s", id_0, id_0_repeated)
	}
	log.Printf("Run IDs: %s, %s\n", id_0, id_1)
}

func TestGetRunHistoryPlayers(t *testing.T) {
	filePath := "example.sav"
	fdata, err := readFile(filePath)
	if err != nil {
		t.Errorf("Could not load testing file: %s\n", err)
	}
	jsonStr, err := ConvertUesaveToJSON(fdata)
	if err != nil {
		t.Fatalf("Failed to parse file: %s\n", err)
	}

	runHistoryData := GetRunHistoryEntries(jsonStr)
	if runHistoryData == nil {
		t.Errorf("Expected run history data, but got an empty string")
	}
	if len(runHistoryData) != 8 {
		t.Errorf("Expected 8 Runs, Got %d\n", len(runHistoryData))
	}
	players0 := GetRunPlayers(runHistoryData[0])
	numUpgrades := []int{19, 22, 17, 20}
	for i, player := range players0 {
		upgrades := GetRunUpgrades(player, GetRunID(runHistoryData[0]))
		if len(upgrades) != numUpgrades[i] {
			t.Errorf("Expected %d upgrades for player %d, but got %d\n", numUpgrades[i], i, len(upgrades))
		}
	}
	items, item_upgrades := GetRunItems(players0[0], GetRunID(runHistoryData[0]))
	log.Println("Items:")
	for _, item := range items {
		log.Printf("Item ID: %s -- Ref: %s\n", item.ItemId, item.Reference)
	}
	log.Println("Upgrades:")
	for _, upgrade := range item_upgrades {
		log.Printf("Upgrade ID: %s, Quantity: %d, Ref: %s\n", upgrade.UpgradeId, upgrade.Quantity, upgrade.Reference)
	}
}
