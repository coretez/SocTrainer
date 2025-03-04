package config

import (
	"context"
	"log"

	firebase "firebase.google.com/go"
	"google.golang.org/api/option"
)

// FirebaseApp is a global variable for Firebase instance
var FirebaseApp *firebase.App

// InitFirebase initializes Firebase
func InitFirebase() {
	opt := option.WithCredentialsFile("config/firebase-adminsdk.json") // ✅ Ensure correct path
	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		log.Fatalf("🔥 Error initializing Firebase: %v", err)
	}

	FirebaseApp = app
	log.Println("✅ Firebase initialized successfully")
}
