package main

import (
	"avito_test_case/config"
	"avito_test_case/internal/app"
	"log"
)

func main() {
	// Configuration
	cfg, err := config.New()
	if err != nil {
		log.Fatalf("Configuration error: %s", err)
	}

	// Run
	app.Run(cfg)
}
