package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/vivienne-curewitz/rogue_core_stats/db"
	"github.com/vivienne-curewitz/rogue_core_stats/types"
)

func UploadHandlerFactory(status map[uuid.UUID]types.UploadStatus, dataPipe chan types.SaveDataTask) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Handle the upload logic here
		requestID := uuid.New()
		status[requestID] = types.UploadStatusPending
		data, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Could not read request body", http.StatusBadRequest)
			return
		}
		dataPipe <- types.SaveDataTask{Data: data, ID: requestID}
		w.WriteHeader(http.StatusAccepted)
		bts, err := w.Write([]byte(requestID.String()))
		if err != nil {
			log.Printf("Failed to write response to upload: %v\n", err)
		}
		if bts == 0 {
			log.Println("Wrote 0 bytes to response for request ID")
		}
	}
}

func UploadStatusHandlerFactory(status map[uuid.UUID]types.UploadStatus) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Handle the status check logic here
		requestIDStr := r.URL.Query().Get("id")
		if requestIDStr == "" {
			http.Error(w, "Missing id parameter", http.StatusBadRequest)
			return
		}
		requestID, err := uuid.Parse(requestIDStr)
		if err != nil {
			http.Error(w, "Invalid id parameter", http.StatusBadRequest)
			return
		}
		if uploadStatus, exists := status[requestID]; exists {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(uploadStatus.String()))
		} else {
			http.Error(w, "Upload not found", http.StatusNotFound)
		}
	}
}

func OverviewHandlerFactory() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Handle the game overview logic here
		playerID := r.URL.Query().Get("player_id")
		if playerID == "" {
			http.Error(w, "Missing player_id parameter", http.StatusBadRequest)
			return
		}
		// get db data here
		overviews, err := db.GetPlayerOverview(r.Context(), playerID)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to fetch overview: %v", err), http.StatusInternalServerError)
			return
		}

		resJSON, err := json.Marshal(overviews)
		if err != nil {
			http.Error(w, "Failed to marshal response", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(resJSON)
	}
}
