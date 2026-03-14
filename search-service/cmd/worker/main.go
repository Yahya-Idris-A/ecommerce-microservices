package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
	"github.com/segmentio/kafka-go"
)

type DebeziumMessage struct {
	Payload struct {
		After map[string]interface{} `json:"after"`
	} `json:"payload"`
}

func main() {
	// 1. Initialize Elasticsearch Client
	es, err := elasticsearch.NewDefaultClient()
	if err != nil {
		log.Fatalf("[FATAL] Error creating the Elasticsearch client: %s", err)
	}

	// 2. Initialize Kafka Reader
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:     []string{"localhost:9092"},
		GroupTopics: []string{"ecommerce.public.products", "ecommerce.public.merchants"},
		GroupID:     "search-service-group",
		MinBytes:    10e3,
		MaxBytes:    10e6,
	})

	fmt.Println("[INFO] Search Service Worker is running...")
	fmt.Println("[INFO] Connected to Elasticsearch at http://localhost:9200")
	fmt.Println("[INFO] Listening to topics: products & merchants")
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

		if msg.Payload.After == nil {
			continue
		}

		data := msg.Payload.After

		// Ambil ID dari map untuk DocumentID di ES
		docID, ok := data["id"].(string)
		if !ok {
			continue // Skip kalau tidak ada field ID
		}

		docBytes, err := json.Marshal(data)
		if err != nil {
			continue
		}

		// Tentukan target Index berdasarkan nama topik Kafka
		var targetIndex string
		if strings.Contains(m.Topic, "products") {
			targetIndex = "products"
		} else if strings.Contains(m.Topic, "merchants") {
			targetIndex = "merchants"
		} else {
			continue
		}

		req := esapi.IndexRequest{
			Index:      targetIndex,
			DocumentID: docID,
			Body:       bytes.NewReader(docBytes),
			Refresh:    "true",
		}

		res, err := req.Do(context.Background(), es)
		if err != nil {
			log.Printf("[ERROR] Failed to index document: %s\n", err)
			continue
		}
		res.Body.Close()

		if !res.IsError() {
			fmt.Printf("[SUCCESS] Indexed to %s | ID: %s\n", targetIndex, docID)
		}
	}

	if err := r.Close(); err != nil {
		log.Fatal("[FATAL] Failed to close reader:", err)
	}
}
