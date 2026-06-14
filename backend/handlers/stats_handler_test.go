package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/vivienne-curewitz/rogue_core_stats/db"
	"github.com/vivienne-curewitz/rogue_core_stats/types"
)

func TestStatsHandler(t *testing.T) {
	ctx := context.Background()
	// start db
	dberr := db.LoadConfig()
	if dberr != nil {
		t.Fatalf("Error with db init: %v\n", dberr)
	}
	dberr = db.InitDB(ctx)
	if dberr != nil {
		t.Fatalf("Error with db init: %v\n", dberr)
	}
	// clear db for testing
	dberr = db.DebugDropTables(ctx)
	if dberr != nil {
		t.Fatalf("Error with db init: %v\n", dberr)
	}
	dberr = db.InitDB(ctx)
	if dberr != nil {
		t.Fatalf("Error with db init: %v\n", dberr)
	}
	// start stats handler
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go StartHandlers(wg)
	// send upload
	data, err := os.ReadFile("../testdata/example.sav")
	if err != nil {
		t.Fatalf("Failed to read test data: %v", err)
	}
	dataReader := bytes.NewReader(data)
	resp, err := http.Post("http://localhost:8080/saveUpload", "application/octect-stream", dataReader)
	if err != nil {
		t.Fatalf("Failed to send upload request: %v", err)
	}
	if resp.StatusCode != http.StatusAccepted {
		t.Fatalf("Expected status code 200, got %d", resp.StatusCode)
	}
	defer resp.Body.Close()
	sz, err := io.ReadAll(resp.Body) // read the UUID from the response body
	if err != nil {
		t.Fatalf("Failed to read request ID from response: %v\n", err)
	}
	id, err := uuid.Parse(string(sz))
	if err != nil {
		t.Fatalf("Failed to read request ID from response: %v\n", err)
	}
	log.Printf("Request ID: %s\n", id)
	// query until result is available
	reqString := fmt.Sprintf("http://localhost:8080/saveUploadStatus?id=%s", id.String())
	for {
		resp, err := http.Get(reqString)
		if err != nil {
			t.Fatalf("Failed to send status request: %v", err)
		}
		defer resp.Body.Close()
		st, err := io.ReadAll(resp.Body)
		if err != nil {
			t.Fatalf("Status handler body error: %v", err)
		}
		status := string(st)
		if status == types.UploadStatusFailed.String() {
			t.Fatalf("Status returned upload failure")
		} else if status == types.UploadStatusCompleted.String() {
			break
		}
		time.Sleep(time.Millisecond * 50)
	}
	// get overview data
	ovReqStr := fmt.Sprintf("http://localhost:8080/overview?player_id=Danger")
	overviewResp, err := http.Get(ovReqStr)
	if err != nil {
		t.Fatalf("Overview Request Failed: %v\n", err)
	}
	// decode into types.RunOverview
	decoder := json.NewDecoder(overviewResp.Body)
	var overViews []types.RunOverview
	var overview types.RunOverview
	err = decoder.Decode(&overViews)
	if err != nil {
		t.Fatalf("Failed to read overview response body: %v\n", err)
	}
	if err != nil {
		t.Fatalf("Failed to unmarshal overview data: %v\n", err)
	}
	overview = overViews[0]
	if overview.PlayerId != "Danger" {
		t.Fatalf("Expected player ID 'Danger', got '%s'", overview.PlayerId)
	}
}
