package main

import (
	"context"
	"fmt"
	"log"

	firebase "firebase.google.com/go"
	"google.golang.org/api/option"
)

func main() {
	// Load Firebase Admin SDK credentials
	opt := option.WithCredentialsFile("config/firebase-adminsdk.json")
	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		log.Fatalf("❌ Failed to initialize Firebase: %v", err)
	}

	// Get Auth client
	client, err := app.Auth(context.Background())
	if err != nil {
		log.Fatalf("❌ Failed to get Auth client: %v", err)
	}

	// Test authentication with an example UID (replace with a real user ID)
	testUID := "AuykIPwKDWX0mP1KRG3WtXM57Ny2"
	userRecord, err := client.GetUser(context.Background(), testUID)
	if err != nil {
		fmt.Println("⚠️ Test failed (user not found):", err)
	} else {
		fmt.Printf("✅ Firebase Auth working! User: %v\n", userRecord)
	}
}
