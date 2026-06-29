package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/vivienne-curewitz/rogue_core_stats/db"
)

func HandleItemOverview(w http.ResponseWriter, r *http.Request) {
	// Get the query parameters
	query := r.URL.Query()
	playerID := query.Get("PlayerId")
	runID := query.Get("RunId")
	if runID == "" || playerID == "" {
		http.Error(w, "Missing RunId or PlayerId parameter", http.StatusBadRequest)
		return
	}

	// from db
	ctx := r.Context()
	items, err := db.GetAssetsByRunIDPlayerID(ctx, runID, playerID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to fetch items: %v", err), http.StatusInternalServerError)
		return
	}
	// encode to JSON and return
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(items)
}

func HandleGetUpgrades(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	playerID := query.Get("PlayerId")
	runID := query.Get("RunId")
	if runID == "" || playerID == "" {
		http.Error(w, "Missing Run or Plaer ID", http.StatusBadRequest)
		return
	}
	// get from db
	ctx := r.Context()
	upgrades, err := db.GetUpgradesByRunIDPlayerID(ctx, runID, playerID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to retreive items from db: %v", err), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(upgrades)
}
