package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"sync"
)

var (
	latestImagePath = "latest.jpg"
	mu              sync.Mutex
)

func uploadHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("📥 Upload received: %s %s", r.Method, r.URL.Path)

	if r.Method != http.MethodPut && r.Method != http.MethodPost {
		http.Error(w, "Only PUT/POST allowed", http.StatusMethodNotAllowed)
		return
	}

	// Lock and save body as file
	mu.Lock()
	defer mu.Unlock()

	f, err := os.Create(latestImagePath)
	if err != nil {
		http.Error(w, "Failed to save image", http.StatusInternalServerError)
		log.Println("❌ Create error:", err)
		return
	}
	defer f.Close()

	_, err = io.Copy(f, r.Body)
	if err != nil {
		http.Error(w, "Failed to write file", http.StatusInternalServerError)
		log.Println("❌ Copy error:", err)
		return
	}

	log.Println("✅ Image saved")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "Image uploaded")
}

func htmlHandler(w http.ResponseWriter, r *http.Request) {
	html := `<!DOCTYPE html>
<html>
<head><title>Latest Upload</title></head>
<body>
<h1>Latest Uploaded Image</h1>
<img src="/latest.jpg" alt="Latest" style="max-width: 100%%; max-height: 500px;">
</body>
</html>`
	w.Header().Set("Content-Type", "text/html")
	fmt.Fprint(w, html)
}

func imageHandler(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	defer mu.Unlock()

	http.ServeFile(w, r, latestImagePath)
}

func main() {
	http.HandleFunc("/", htmlHandler)
	http.HandleFunc("/upload", uploadHandler)
	http.HandleFunc("/latest.jpg", imageHandler)

	log.Println("🚀 Server running on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
