package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"backend/config"
	"backend/replay"

	"google.golang.org/api/iterator"
)

// Global channel for progress updates
var progressChan chan replay.ReplayProgress

// FetchScenarioFile retrieves the scenario file URL from Firestore
func FetchScenarioFile(scenarioName string) (string, error) {
	ctx := context.Background()
	client, err := config.FirebaseApp.Firestore(ctx) // ✅ Assign both values
	if err != nil {
		return "", fmt.Errorf("error initializing Firestore client: %v", err)
	}
	defer client.Close()

	scenariosRef := client.Collection("scenarios")
	query := scenariosRef.Where("name", "==", scenarioName).Documents(ctx)

	for {
		doc, err := query.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return "", fmt.Errorf("error fetching scenario: %v", err)
		}

		scenarioData := doc.Data()
		if fileURL, ok := scenarioData["stream"].(string); ok {
			return fileURL, nil
		}
	}

	return "", fmt.Errorf("scenario not found")
}

// DownloadFile downloads a file from Firebase Storage
func DownloadFile(fileURL string) (string, error) {
	resp, err := http.Get(fileURL)
	if err != nil {
		return "", fmt.Errorf("failed to download file: %v", err)
	}
	defer resp.Body.Close()

	tempFilePath := filepath.Join(os.TempDir(), "scenario.json")
	file, err := os.Create(tempFilePath)
	if err != nil {
		return "", fmt.Errorf("failed to create temp file: %v", err)
	}
	defer file.Close()

	_, err = io.Copy(file, resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to save scenario file: %v", err)
	}

	return tempFilePath, nil
}

// ReplayHandler starts a replay using Firebase Storage URL
func ReplayHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		HECToken     string `json:"hec_token"`
		HECURL       string `json:"hec_url"`
		ScenarioName string `json:"scenario_name"`
	}

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Fetch scenario file URL from Firestore
	fileURL, err := FetchScenarioFile(req.ScenarioName)
	if err != nil {
		http.Error(w, "Scenario not found in Firestore", http.StatusNotFound)
		return
	}

	// Download scenario file locally
	localFilePath, err := DownloadFile(fileURL)
	if err != nil {
		http.Error(w, "Failed to download scenario file", http.StatusInternalServerError)
		return
	}

	// Start a new progress channel
	progressChan = make(chan replay.ReplayProgress)

	// Start replay in a goroutine
	go replay.ReplayRecords(localFilePath, req.HECURL, req.HECToken, progressChan)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Replay started successfully"})
}

// ProgressHandler streams replay progress using SSE
func ProgressHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	flusher, _ := w.(http.Flusher)
	for progress := range progressChan {
		progressJSON, _ := json.Marshal(progress)
		fmt.Fprintf(w, "data: %s\n\n", progressJSON)
		flusher.Flush()
	}
}