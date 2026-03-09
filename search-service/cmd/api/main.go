package main

import (
	"fmt"
	"log"

	delivery "github.com/Yahya-idris-A/ecommerce-microservices/search-service/internal/delivery/http"
	"github.com/Yahya-idris-A/ecommerce-microservices/search-service/internal/repository/elasticsearch"
	"github.com/Yahya-idris-A/ecommerce-microservices/search-service/internal/usecase"
	es "github.com/elastic/go-elasticsearch/v8"
	"github.com/gin-gonic/gin"
)

func main() {
	// 1. Setup Infrastructure
	esClient, err := es.NewDefaultClient()
	if err != nil {
		log.Fatalf("[FATAL] Error creating the Elasticsearch client: %s", err)
	}

	// 2. Setup Repository
	searchRepo := elasticsearch.NewProductSearchRepository(esClient)

	// 3. Setup Usecase
	searchUsecase := usecase.NewProductSearchUsecase(searchRepo)

	// 4. Setup Delivery (Gin Router)
	router := gin.Default()
	delivery.NewSearchHandler(router, searchUsecase)

	// 5. Start Server
	fmt.Println("[INFO] Search API is running on port 8081 with Clean Architecture")
	if err := router.Run(":8081"); err != nil {
		log.Fatalf("[FATAL] Failed to start server: %v", err)
	}
}
