package main

import (
	"log"

	"github.com/andreychh/coopera-backend/internal/app"
)

func main() {
	if err := app.Start(); err != nil {
		log.Fatalf("Failed to start the application: %v", err)
	}
}
