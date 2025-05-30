package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	minio "github.com/minio/minio-go/v7"
	credentials "github.com/minio/minio-go/v7/pkg/credentials"
)

var (
	minioAccessKeyID     string
	minioSecretAccessKey string
	minioEndpoint        string
	minioBucketName      string
	useSSL               = false
	minioClient          *minio.Client
	err                  error
)

func init() {
	minioAccessKeyID = loadEnvVar("MINIO_ACCESS_KEY_ID")
	minioSecretAccessKey = loadEnvVar("MINIO_SECRET_ACCESS_KEY")
	minioEndpoint = loadEnvVar("MINIO_ENDPOINT")
	minioBucketName = loadEnvVar("MINIO_BUCKET_NAME")
}

func loadEnvVar(key string) string {
	value, ok := os.LookupEnv(key)
	if !ok {
		log.Fatalf("env var %s not set. Exiting.", key)
	}
	return value
}

func uploadHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("📥 Upload received: %s %s", r.Method, r.URL.Path)

	if r.Method != http.MethodPut && r.Method != http.MethodPost {
		http.Error(w, "Only PUT/POST allowed", http.StatusMethodNotAllowed)
		return
	}

	minioClient, err = minio.New(minioEndpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(minioAccessKeyID, minioSecretAccessKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		log.Fatalf("❌ Failed to initialize MinIO client: %v", err)
	}
	// Set object name based on time (or customize)
	objectName := fmt.Sprintf("upload-%d", time.Now().UnixNano())

	// Upload to MinIO
	info, err := minioClient.PutObject(
		context.Background(),
		minioBucketName,
		objectName,
		r.Body,
		r.ContentLength,
		minio.PutObjectOptions{ContentType: "application/octet-stream"},
	)
	if err != nil {
		log.Printf("❌ Upload failed: %v", err)
		http.Error(w, "Upload failed", http.StatusInternalServerError)
		return
	}

	log.Printf("✅ Uploaded object %s (%d bytes)", info.Key, info.Size)

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "OK: %s (%d bytes uploaded)\n", info.Key, info.Size)
}

func main() {
	http.HandleFunc("/upload", uploadHandler)

	log.Println("🚀 Server running on :80")
	log.Fatal(http.ListenAndServe(":80", nil))
}
