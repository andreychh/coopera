package main

import (
	"log"

	"github.com/andreychh/coopera-bot/internal/app"
)

func main() {
	err := app.Start()
	if err != nil {
		log.Fatalf("starting app: %v\n", err)
	}
}
