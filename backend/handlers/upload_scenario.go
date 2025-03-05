package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"cloud.google.com/go/storage"
)

// UploadResponse represents the response after uploading a file.
type UploadResponse struct {
	FileURL string `json:"file_url"`
}

// UploadScenarioHandler handles file uploads to Firebase Storage.
func UploadScenarioHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("📥 Incoming upload request...")

	// Ensure it's a POST request
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	// Parse the multipart form
	err := r.ParseMultipartForm(10 << 20) // 10MB max file size
	if err != nil {
		log.Printf("❌ Failed to parse form: %v", err)
		http.Error(w, "Failed to parse form", http.StatusBadRequest)
		return
	}

	// Retrieve the file
	file, handler, err := r.FormFile("file")
	if err != nil {
		log.Printf("❌ Failed to get uploaded file: %v", err)
		http.Error(w, "Failed to get uploaded file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	fileName := handler.Filename
	log.Printf("📂 File received: %s (%d bytes)", fileName, handler.Size)

	// Upload to Firebase Storage using the improved function
	uploadedURL, err := UploadFileToFirebase("replaydata-385e9.firebasestorage.app", fileName, file)
	if err != nil {
		log.Printf("❌ Upload to Firebase failed: %v", err)
		http.Error(w, fmt.Sprintf("Failed to upload file: %v", err), http.StatusInternalServerError)
		return
	}

	// Respond with the uploaded file URL
	response := UploadResponse{FileURL: uploadedURL}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
	log.Println("✅ File uploaded successfully:", uploadedURL)
}

func UploadFileToFirebase(bucketName, fileName string, file io.Reader) (string, error) {
	ctx := context.Background()

	client, err := storage.NewClient(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to create storage client: %v", err)
	}
	defer client.Close()

	bucket := client.Bucket(bucketName) //The bucketName is just the name, without gs://
	obj := bucket.Object(fmt.Sprintf("scenarios/%s", fileName))
	wc := obj.NewWriter(ctx)

	n, err := io.Copy(wc, file)
	if err != nil {
		return "", fmt.Errorf("failed to write file: %w. Wrote %d bytes", err, n)
	}

	if err := wc.Close(); err != nil {
		return "", fmt.Errorf("failed to close writer: %v", err)
	}

	// Construct the URL correctly using the bucketName without the gs:// prefix.
	url := fmt.Sprintf("https://storage.googleapis.com/%s/%s", bucketName, fileName)
	return url, nil
}
