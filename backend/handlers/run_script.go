package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os/exec"
)

// RunScriptHandler triggers a replay process
func RunScriptHandler(w http.ResponseWriter, r *http.Request) {
	// Decode request
	var request struct {
		HECToken     string `json:"hec_token"`
		HECURL       string `json:"hec_url"`
		ScenarioFile string `json:"scenario_file"`
	}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Run replay as a background process
	go func() {
		cmd := exec.Command("go", "run", "replay.go", request.HECToken, request.HECURL, request.ScenarioFile)
		if err := cmd.Run(); err != nil {
			fmt.Println("Error running script:", err)
		}
	}()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Script execution started"})
}
