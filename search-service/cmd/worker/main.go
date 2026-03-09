package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
	"github.com/segmentio/kafka-go"
)

type Product struct {
	ID          string `json:"id"`
	CategoryID  string `json:"category_id"`
	Name        string `json:"name"`
	Slug        string `json:"slug"`
	Description string `json:"description"`
	Price       string `json:"price"`
	Stock       int32  `json:"stock"`
	MerchantID  string `json:"merchant_id"`
}

type DebeziumPayload struct {
	Op    string   `json:"op"`
	After *Product `json:"after"`
}

type DebeziumMessage struct {
	Payload DebeziumPayload `json:"payload"`
}

func main() {
	// 1. Initialize Elasticsearch Client
	es, err := elasticsearch.NewDefaultClient()
	if err != nil {
		log.Fatalf("[FATAL] Error creating the Elasticsearch client: %s", err)
	}

	// 2. Initialize Kafka Reader
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  []string{"localhost:9092"},
		Topic:    "ecommerce.public.products",
		GroupID:  "search-service-group",
		MinBytes: 10e3,
		MaxBytes: 10e6,
	})

	fmt.Println("[INFO] Search Service Worker is running...")
	fmt.Println("[INFO] Connected to Elasticsearch at http://localhost:9200")
	fmt.Println("[INFO] Listening to Kafka topic: ecommerce.public.products")
	fmt.Println("---------------------------------------------------")

	for {
		m, err := r.ReadMessage(context.Background())
		if err != nil {
			log.Printf("[ERROR] Failed to read message: %v\n", err)
			break
		}

		var msg DebeziumMessage
		if err := json.Unmarshal(m.Value, &msg); err != nil {
			continue
		}

		// Skip if there is no data (e.g., delete operation)
		if msg.Payload.After == nil {
			continue
		}

		product := msg.Payload.After

		// 3. Convert clean product struct back to JSON for Elasticsearch
		docBytes, err := json.Marshal(product)
		if err != nil {
			log.Printf("[ERROR] Failed to marshal product: %v\n", err)
			continue
		}

		// 4. Send data to Elasticsearch Index named "products"
		req := esapi.IndexRequest{
			Index:      "products",
			DocumentID: product.ID, // Use the same ID from PostgreSQL
			Body:       bytes.NewReader(docBytes),
			Refresh:    "true", // Make it searchable immediately
		}

		res, err := req.Do(context.Background(), es)
		if err != nil {
			log.Printf("[ERROR] Failed to index document: %s\n", err)
			continue
		}
		res.Body.Close()

		if res.IsError() {
			log.Printf("[ERROR] Error indexing document ID=%s: %s\n", product.ID, res.String())
		} else {
			fmt.Printf("[SUCCESS] Indexed to Elasticsearch! ID: %s | Name: %s\n", product.ID, product.Name)
		}
	}

	if err := r.Close(); err != nil {
		log.Fatal("[FATAL] Failed to close reader:", err)
	}
}
