package main

import (
	"github.com/andreychh/coopera/internal/app"
	"log"
)

func main() {
	if err := app.Start(); err != nil {
		log.Fatalf("Failed to start the application: %v", err)
	}
}
