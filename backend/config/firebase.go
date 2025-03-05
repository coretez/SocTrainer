package config

import (
	"context"
	"log"
	"os"

	firebase "firebase.google.com/go"
	"google.golang.org/api/option"
)

// FirebaseApp is a global variable for Firebase instance
var FirebaseApp *firebase.App

// InitFirebase initializes Firebase
func InitFirebase() {
	os.Setenv("GOOGLE_APPLICATION_CREDENTIALS", "config/firebase-adminsdk.json")
	log.Printf("GOOGLE_APPLICATION_CREDENTIALS: %s", os.Getenv("GOOGLE_APPLICATION_CREDENTIALS"))

	opt := option.WithCredentialsFile("config/firebase-adminsdk.json") // Ensure correct path
	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		log.Fatalf("🔥 Error initializing Firebase: %v", err)
	}

	FirebaseApp = app
	log.Println("✅ Firebase initialized successfully")
}

// GetFirebaseApp returns the FirebaseApp instance
func GetFirebaseApp() *firebase.App {
	if FirebaseApp == nil {
		log.Println("⚠️ FirebaseApp is not initialized. Initializing now...")
		InitFirebase()
	}
	return FirebaseApp
}
