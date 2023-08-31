package main

import (
	"avito_test_case/config"
	"avito_test_case/internal/app"
	"log"
)

// @title Denis Saltykov's Solution
// @version 1.0
// @description API Server for User Segmentation

// @host localhost:8080
// @BasePath /

func main() {
	// Configuration
	cfg, err := config.New()
	if err != nil {
		log.Fatalf("Configuration error: %s", err)
	}

	// Run
	app.Run(cfg)
}
