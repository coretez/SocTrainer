package handlers

import (
	"fmt"
	"net/http"
)

// FetchScenarioHandler - Sample placeholder handler
func FetchScenarioHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "✅ FetchScenarioHandler is working!")
}
