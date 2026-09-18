package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/dstepan499-dev/url-shortener-api/internal/handler"
	"github.com/dstepan499-dev/url-shortener-api/internal/repository"
	"github.com/dstepan499-dev/url-shortener-api/internal/service"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
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

	shortenerService := service.NewShortenerService(repo)
	baseURL := fmt.Sprintf("http://localhost:%s", port)
	h := handler.NewHandler(shortenerService, baseURL)

	server := &http.Server{
		Addr:    ":" + port,
		Handler: h.InitRoutes(),
	}

	log.Printf("Server is running on port %s...", port)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Server failed to start: %v", err)
	}
}
