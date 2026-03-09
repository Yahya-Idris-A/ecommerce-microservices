package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	delivery "github.com/Yahya-idris-A/ecommerce-microservices/user-service/internal/delivery/http"
	"github.com/Yahya-idris-A/ecommerce-microservices/user-service/internal/repository/postgres"
	"github.com/Yahya-idris-A/ecommerce-microservices/user-service/internal/usecase"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	// 1. Load Environment Variables
	if err := godotenv.Load(); err != nil {
		log.Println("[WARNING] No .env file found. Falling back to system environment variables.")
	}

	dbURL := os.Getenv("DB_URL")
	jwtSecret := os.Getenv("JWT_SECRET")

	if dbURL == "" || jwtSecret == "" {
		log.Fatal("[FATAL] DB_URL and JWT_SECRET must be set in environment")
	}

	// 2. Initialize Database Connection
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("[FATAL] Failed to open database connection: %v", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatalf("[FATAL] Failed to ping database: %v", err)
	}
	fmt.Println("[INFO] Successfully connected to PostgreSQL User Database")

	// 3. Setup Clean Architecture Layers
	userRepo := postgres.NewUserRepository(db)
	userUsecase := usecase.NewUserUsecase(userRepo, jwtSecret)

	// 4. Setup Gin Router
	router := gin.Default()

	// Ini adalah pemanggilan fungsi yang benar
	delivery.NewUserHandler(router, userUsecase)

	// 5. Start the Server on Port 8082
	fmt.Println("[INFO] User Service API is running on port 8082")
	if err := router.Run(":8082"); err != nil {
		log.Fatalf("[FATAL] Failed to start server: %v", err)
	}
}
