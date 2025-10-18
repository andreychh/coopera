package app

import (
	"encoding/json"
	"fmt"
	"github.com/andreychh/coopera/internal/adapter/controller/telegram_api"
	"github.com/andreychh/coopera/internal/adapter/repository/postgres"
	"github.com/andreychh/coopera/pkg/logger"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
	"strconv"
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

func Start() error {
	log.Println("Service starting...")
	log.Printf("Build info: version=%s, commit=%s", version, commit)

	err := godotenv.Load("config/dev/.env")

	if err != nil {
		return fmt.Errorf("error loading .env file")
	}

	logLevel, err := strconv.Atoi(os.Getenv("LOG_LEVEL"))
	if err != nil {
		logLevel = logger.INFO
	}
	logService := logger.NewLogger(logLevel)

	botToken := os.Getenv("TELEGRAM_BOT_TOKEN")
	if botToken == "" {
		logService.Fatal("Bot token not found in environment variables")
	}

	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		logService.Fatal("Database URL not found in environment variables")
	}

	_, err = postgres.NewDB(dsn)
	if err != nil {
		return err
	}

	tgContr, err := telegram_api.NewTelegramController(logService, botToken)
	if err != nil {
		return err
	}

	go func() {
		err = tgContr.Start()
		if err != nil {
			log.Printf("err in bot")
		}
	}()

	http.HandleFunc("/status", statusHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s...", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatalf("FATAL: could not start server: %v", err)
	}

	return nil
}
