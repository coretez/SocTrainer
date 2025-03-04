package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()

	// API Routes
	r.HandleFunc("/", HomeHandler).Methods("GET")
	r.HandleFunc("/api/run", RunScriptHandler).Methods("POST")
	r.HandleFunc("/api/logs", StreamLogsHandler).Methods("GET")
	r.HandleFunc("/api/upload-scenario", UploadScenarioHandler).Methods("POST")

	// Enable CORS
	handler := cors.Default().Handler(r)

	fmt.Println("🚀 Server started on port 5000")
	log.Fatal(http.ListenAndServe(":5000", handler))
}