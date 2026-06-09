package saveparser

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"log"
	"slices"
	"strings"

	"github.com/google/uuid"
	"github.com/tidwall/gjson"
	"github.com/vivienne-curewitz/rogue_core_stats/db"
	"github.com/vivienne-curewitz/rogue_core_stats/types"
)

// saveData is a json string
func GetRunHistoryEntries(saveData string) []string {
	runs := gjson.Get(saveData, "root.properties.RunHistory_0.Entries_0")
	if runs.Exists() {
		retval := make([]string, 0)
		runs.ForEach(func(key, value gjson.Result) bool {
			retval = append(retval, value.String())
			return true
		})
		return retval
	} else {
		return nil
	}
}

// sort the player names alphanumerically, append the runtime and seed, and hash with sha256
func GetRunID(runString string) string {
	nameResults := gjson.Get(runString, "Characters_0.BuildData_0.#.PlayerName_0").Array()
	timestamp := gjson.Get(runString, "SharedInfo_0.MissionTime_0").String()
	seed := gjson.Get(runString, "RunInfo_0.RunSeed_0").String()
	names := make([]string, len(nameResults))
	for i, nr := range nameResults {
		names[i] = nr.String()
	}
	// sort the names
	slices.Sort(names)
	// concatenate the names, timestamp, and seed
	concat := strings.Join(names, "") + timestamp + seed
	// hash with sha256
	hash := sha256.Sum256([]byte(concat))
	return hex.EncodeToString(hash[:])
}

func GetRunStatus(runString string) bool {
	return gjson.Get(runString, "SharedInfo_0.MissionSuccess_0").Bool()
}

func GetRunPlayers(runString string) []string {
	playerData := gjson.Get(runString, "Characters_0.BuildData_0").Array()
	retPd := make([]string, len(playerData))
	for i, pd := range playerData {
		retPd[i] = pd.String()
	}
	return retPd
}

func GetRunUpgrades(playerString string, runID string) []types.Upgrade {
	name := gjson.Get(playerString, "PlayerName_0")
	upgradesRet := gjson.Get(playerString, "UnlockRecords_0.#(SelectedSlot_0.SlotIndex_0==-1)#").Array()
	retUps := make([]types.Upgrade, len(upgradesRet))
	for i, ur := range upgradesRet {
		itemID := ur.Get("Unlock_0")
		quantity := ur.Get("count_0")
		retUps[i] = types.Upgrade{
			RunId:     runID,
			PlayerId:  name.String(),
			UpgradeId: itemID.String(),
			Quantity:  int(quantity.Int()),
		}
	}
	return retUps
}

func GetRunItems(playerString string, runID string) ([]types.Item, []types.Upgrade) {
	name := gjson.Get(playerString, "PlayerName_0")
	itemsRet := gjson.Get(playerString, "UnlockRecords_0.#(SelectedSlot_0.SlotIndex_0!=-1)#").Array()
	retItems := make([]types.Item, len(itemsRet))
	retUps := make([]types.Upgrade, 0)
	for i, ur := range itemsRet {
		itemID := ur.Get("Unlock_0")
		iref := uuid.New().String()
		retItems[i] = types.Item{
			RunId:     runID,
			PlayerId:  name.String(),
			ItemId:    itemID.String(),
			Reference: iref,
		}
		itemUpgrades := ur.Get("Attributes_0").Array()
		for _, iu := range itemUpgrades {
			upgradeID := iu.Get("ID_0")
			quantity := iu.Get("count_0")
			retUps = append(retUps, types.Upgrade{
				RunId:     runID,
				PlayerId:  name.String(),
				UpgradeId: upgradeID.String(),
				Quantity:  int(quantity.Int()),
				Reference: iref,
			})
		}
	}
	return retItems, retUps
}

func consolidateUpgrades(upgrades []types.Upgrade) []types.Upgrade {
	m := make(map[string]types.Upgrade)
	for _, u := range upgrades {
		key := u.RunId + "|" + u.PlayerId + "|" + u.UpgradeId + "|" + u.Reference
		if existing, ok := m[key]; ok {
			existing.Quantity += u.Quantity
			m[key] = existing
		} else {
			m[key] = u
		}
	}
	res := make([]types.Upgrade, 0, len(m))
	for _, v := range m {
		res = append(res, v)
	}
	return res
}

func ExtractRunData(runString string) (bool, error) {
	ctx := context.Background()
	runId := GetRunID(runString)
	exists, err := db.RunExists(ctx, runId)
	if err != nil {
		log.Printf("Error checking if run exists: %v", err)
		return false, err
	}
	if exists {
		return false, nil
	}
	runStatus := GetRunStatus(runString)
	err = db.WriteRunStatus(ctx, types.RunStatus{
		RunId:  runId,
		Status: runStatus,
	})
	if err != nil {
		log.Printf("Error inserting run status: %v", err)
		return false, err
	}
	// check if runId already exists in the database
	players := GetRunPlayers(runString)
	for _, player := range players {
		upgrades := GetRunUpgrades(player, runId)
		items, itemUpgrades := GetRunItems(player, runId)
		upgrades = slices.Concat(upgrades, itemUpgrades)
		// consolidate upgrades to handle duplicates (e.g. same upgrade on same item)
		upgrades = consolidateUpgrades(upgrades)
		// write items
		db.BatchWriteItems(ctx, items)
		// write upgrades
		db.BatchWriteUpgrades(ctx, upgrades)
	}
	return true, nil
}
