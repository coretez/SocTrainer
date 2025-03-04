package main

import (
	"fmt"
	"log"
	"net/http"

	"backend/config"
	"backend/handlers"

	ghandlers "github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

func main() {
	fmt.Println("🚀 Initializing Firebase...")
	config.InitFirebase()

	router := mux.NewRouter()

	// Register API routes
	router.HandleFunc("/api/replay", handlers.ReplayHandler).Methods("POST")
	router.HandleFunc("/api/replay/progress", handlers.ProgressHandler).Methods("GET")

	// CORS Middleware
	corsHandler := ghandlers.CORS(
		ghandlers.AllowedOrigins([]string{"http://localhost:8888"}),
		ghandlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}),
		ghandlers.AllowedHeaders([]string{"Content-Type", "Authorization"}),
	)

	port := ":8080"
	fmt.Println("🚀 Server running on http://localhost" + port)
	log.Fatal(http.ListenAndServe(port, corsHandler(router)))
}