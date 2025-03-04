package handlers

import (
	"encoding/json"
	"net/http"
)

// StreamLogsHandler handles streaming logs
func StreamLogsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Log streaming started"})
}

// StopLogsHandler stops log streaming
func StopLogsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Log streaming stopped"})
}
