package main

import (
	"log"
	"os"

	"github.com/qmni/swe.workshop/internal/app"
	"github.com/qmni/swe.workshop/internal/database"
	"github.com/qmni/swe.workshop/internal/middleware"
)

func main() {
	cfg := database.ConfigFromEnv()
	db, err := database.Open(cfg)
	if err != nil {
		log.Fatalf("connect database: %v", err)
	}

	if err := database.Migrate(db); err != nil {
		log.Fatalf("migrate database: %v", err)
	}

	if err := middleware.InitKeycloak(); err != nil {
		log.Fatalf("init keycloak: %v", err)
	}

	server := app.New(db)
	addr := ":" + env("PORT", "8080")
	if err := server.Listen(addr); err != nil {
		log.Fatalf("listen %s: %v", addr, err)
	}
}

func env(key string, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}
