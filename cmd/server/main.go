package main

import (
	"context"
	"log"
	"os"

	"github.com/clerk/clerk-sdk-go/v2"
	"github.com/joho/godotenv"

	"mova-backend/internal/database"
	"mova-backend/internal/repository"
	"mova-backend/internal/router"
	"mova-backend/internal/service"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("DATABASE_URL is required")
	}

	port := "8080"
	if os.Getenv("PORT") != "" {
		port = os.Getenv("PORT")
	}

	clerkKey := os.Getenv("CLERK_SECRET_KEY")
	if clerkKey == "" {
		log.Fatal("CLERK_SECRET_KEY is required")
	}
	clerk.SetKey(clerkKey)

	clerkWebhookSecret := os.Getenv("CLERK_WEBHOOK_SECRET")
	if clerkWebhookSecret == "" {
		log.Fatal("CLERK_WEBHOOK_SECRET is required")
	}

	ctx := context.Background()
	pool, err := database.NewPool(ctx, dbURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer pool.Close()

	log.Println("Connected to database")

	// Initialize layers
	queries := database.New(pool)
	userRepo := repository.NewUserRepository(queries)
	userService := service.NewUserService(userRepo)

	r, err := router.Setup(userService, clerkWebhookSecret)
	if err != nil {
		log.Fatalf("Failed to setup router: %v", err)
	}

	log.Println("Starting server on :" + port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
