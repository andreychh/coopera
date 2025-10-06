package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"
)

// Build information, injected at compile time.
var (
	version = "(devel)"
	commit  = "none"
)

// StatusResponse is the structure of our /status endpoint response.
type StatusResponse struct {
	Status    string    `json:"status"`
	Version   string    `json:"version"`
	Commit    string    `json:"commit"`
	Timestamp time.Time `json:"timestamp"`
}

func statusHandler(w http.ResponseWriter, r *http.Request) {
	response := StatusResponse{
		Status:    "ok",
		Timestamp: time.Now().UTC(),
		Version:   version,
		Commit:    commit,
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("ERROR: could not encode status response: %v", err)
	}
}

func main() {
	log.Println("Service starting...")
	log.Printf("Build info: version=%s, commit=%s", version commit)

	http.HandleFunc("/status", statusHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s...", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatalf("FATAL: could not start server: %v", err)
	}
}
