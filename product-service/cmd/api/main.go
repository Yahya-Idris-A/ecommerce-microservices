package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"

	// Ingat untuk selalu menyesuaikan path import ini dengan nama module-mu
	"github.com/Yahya-idris-A/ecommerce-microservices/product-service/internal/config"
	delivery "github.com/Yahya-idris-A/ecommerce-microservices/product-service/internal/delivery/http"
	"github.com/Yahya-idris-A/ecommerce-microservices/product-service/internal/repository"
	"github.com/Yahya-idris-A/ecommerce-microservices/product-service/internal/usecase"
)

func main() {
	// 1. Muat konfigurasi dari file .env
	cfg := config.LoadConfig()

	// 2. Inisialisasi Database menggunakan konfigurasi yang dinamis
	db := config.InitDB(cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName)
	defer db.Close()

	// 3. Setup Dependency Injection
	productRepo := repository.NewPostgresProductRepository(db)
	productUsecase := usecase.NewProductUsecase(productRepo)

	if os.Getenv("APP_ENV") == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	// 4. Inisialisasi Router Gin
	router := gin.Default()

	// Mengamankan proxy (karena kita jalan di lokal, kita nonaktifkan dulu trust proxy-nya)
	// Di production, kamu bisa isi dengan IP Nginx/Load Balancer kamu.
	_ = router.SetTrustedProxies(nil)

	// 5. Mendaftarkan Handler
	delivery.NewProductHandler(router, productUsecase)

	// 6. Jalankan Server HTTP sesuai port di .env
	log.Printf("starting product-service on port %s", cfg.AppPort)
	if err := router.Run(":" + cfg.AppPort); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
