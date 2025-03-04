package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// UploadScenarioHandler handles scenario file uploads
func UploadScenarioHandler(w http.ResponseWriter, r *http.Request) {
	// Check for file in request
	file, _, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Failed to read file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Process file (you can store it, parse it, etc.)
	fmt.Println("✅ Scenario file uploaded")

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Scenario uploaded successfully"})
}
