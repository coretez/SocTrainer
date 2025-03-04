package handlers

import (
	"context"
	"fmt"
	"net/http"

	"backend/config"
)

// Example API that uses Firebase Authentication
func ExampleHandler(w http.ResponseWriter, r *http.Request) {
	// Ensure Firebase is initialized
	if config.FirebaseApp == nil {
		http.Error(w, "Firebase not initialized", http.StatusInternalServerError)
		return
	}

	// Get Firebase Auth client
	client, err := config.FirebaseApp.Auth(context.Background()) // ✅ FIXED
	if err != nil {
		http.Error(w, "Failed to get Firebase Auth client", http.StatusInternalServerError)
		fmt.Println("Firebase Auth error:", err)
		return
	}

	// Use the client to print success message (remove warning)
	_ = client // ✅ This prevents "declared but not used" error

	// ✅ Successfully initialized Firebase Auth
	fmt.Fprintln(w, "✅ Firebase Auth initialized and ready to use")
}
