package main

import (
	"log"
	"os"

	"github.com/dstepan499-dev/url-shortener-api/internal/repository"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found")
	}

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("DATABASE_URL is not set")
	}

	repo, err := repository.New(dbURL)
	if err != nil {
		log.Fatalf("Database connection failed: %v", err)
	}
	defer repo.Close()

	migrationsSQL, err := os.ReadFile("migrations/000001_init.up.sql")
	if err != nil {
		log.Fatalf("Failed to read migration file: %v", err)
	}

	if err := repo.InitSchema(string(migrationsSQL)); err != nil {
		log.Fatalf("Failed to apply migration: %v", err)
	}

	log.Println("Successfully connected to PostgreSQL and initialized schema")
}
