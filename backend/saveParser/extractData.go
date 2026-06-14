package saveparser

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log"
	"slices"
	"strconv"
	"strings"
	"time"

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

func GetRunOverview(runString string, localID int, RunId string, PlayerId string, Cid string) types.RunOverview {
	si := gjson.Get(runString, "SharedInfo_0")
	ri := gjson.Get(runString, "RunInfo_0")
	pSearchStr := fmt.Sprintf("BasicCharacterStats_0.#(PlayerId_0==%d)", localID)
	pi := si.Get(pSearchStr)
	depth := ri.Get("RunDepth_0").String()
	d2 := strings.Replace(depth, "ERunDepth::Depth", "", 1)
	depthInt, err := strconv.Atoi(d2)
	if err != nil {
		log.Printf("Failed to parse depth from string: %s", depth)
	}
	day := si.Get("Day_0").Int()
	month := si.Get("Month_0").Int()
	year := si.Get("Year_0").Int()
	timestamp := time.Date(int(year), time.Month(month), int(day), 0, 0, 0, 0, time.UTC).Unix()
	return types.RunOverview{
		RunId:           RunId,
		PlayerId:        PlayerId,
		Status:          si.Get("MissionSuccess_0").Bool(),
		BossId:          "Unk",
		Depth:           int(depthInt), // need to parse int here
		CharacterId:     Cid,
		PlayerDamage:    float32(pi.Get("TotalCappedDamage_0").Float()),
		OverkillDamage:  float32(pi.Get("TotalOverkillDamage_0").Float()),
		PlayerKills:     int32(pi.Get("TotalKills_0").Int()),
		PlayerDeaths:    int32(pi.Get("TotalDeaths_0").Int()),
		CompletedStages: int32(si.Get("CompletedStages_0").Int()),
		Runtime:         int32(si.Get("MissionTime_0").Int()),
		PlayerRank:      int32(pi.Get("PlayerRank_0").Int()),
		CharacterRank:   int32(pi.Get("CharacterLevel_0").Int()),
		CharacterStars:  int32(pi.Get("Stars_0").Int()),
		MineralsMined:   float32(pi.Get("TotalMineralsMined_0").Float()),
		MaxArmor:        float32(pi.Get("MaxArmor_0").Float()),
		MaxHealth:       float32(pi.Get("MaxHealth_0").Float()),
		HealthRestored:  float32(pi.Get("TotalHealthRestored_0").Float()),
		Timestamp:       timestamp,
	}
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

func GetPlayerIDs(player string) (int, string, string) {
	localID := int(gjson.Get(player, "PlayerId_0").Int())
	playerID := gjson.Get(player, "PlayerName_0").String()
	characterID := gjson.Get(player, "PlayerCharacterId_0").String()
	return localID, playerID, characterID
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
		localId, pid, cid := GetPlayerIDs(player)
		runOverview := GetRunOverview(runString, localId, runId, pid, cid)
		upgrades := GetRunUpgrades(player, runId)
		items, itemUpgrades := GetRunItems(player, runId)
		upgrades = slices.Concat(upgrades, itemUpgrades)
		// consolidate upgrades to handle duplicates (e.g. same upgrade on same item)
		upgrades = consolidateUpgrades(upgrades)
		// write items
		db.BatchWriteItems(ctx, items)
		// write upgrades
		db.BatchWriteUpgrades(ctx, upgrades)
		// write runOverview
		db.BatchWriteRunInfo(ctx, []types.RunOverview{runOverview})
	}
	return true, nil
}

func ReverseSaveRead(task types.SaveDataTask) {
	jsonStr, err := ConvertUesaveToJSON(task.Data)
	if err != nil {
		// statuses -- requests id
		log.Printf("Failed to decode save file: %v\n", err)
	}
	runs := GetRunHistoryEntries(jsonStr)
	for i := len(runs) - 1; i >= 0; i -= 1 {
		ExtractRunData(runs[i])
	}
}

func SaveDataPipe(dataPipe chan types.SaveDataTask, statuses map[uuid.UUID]types.UploadStatus, ctx context.Context) {
MainLoop:
	for {
		select {
		case task := <-dataPipe:
			func() {
				jsonStr, err := ConvertUesaveToJSON(task.Data)
				if err != nil {
					// statuses -- requests id
					log.Printf("Failed to decode save file: %v\n", err)
					statuses[task.ID] = types.UploadStatusFailed
				} else {
					statuses[task.ID] = types.UploadStatusInProgress
					runs := GetRunHistoryEntries(jsonStr)
					for i := len(runs) - 1; i >= 0; i -= 1 {
						ExtractRunData(runs[i])
					}
					statuses[task.ID] = types.UploadStatusCompleted
				}
			}()
		case <-ctx.Done():
			break MainLoop
		}
	}
}
