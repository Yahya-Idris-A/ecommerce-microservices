## How to Start the E-Commerce Microservices Locally

When starting the development environment after a system reboot, follow these exact steps to ensure all services and event streams initialize correctly.

### 1. Start Infrastructure
Open a terminal in the root directory (`ecommerce-microservices`) and start all Docker containers:
```bash
docker compose up -d

# Mengecek status container, pastikan semuanya 'Up'
docker compose ps

# Mendaftarkan konektor Debezium untuk mengawasi tabel products hanya dijalankan setelah mengecek status konektor dan hasilnya tidak running
curl -i -X POST -H "Accept:application/json" -H "Content-Type:application/json" localhost:8083/connectors/ -d '{
  "name": "product-connector",
  "config": {
    "connector.class": "io.debezium.connector.postgresql.PostgresConnector",
    "database.hostname": "postgres",
    "database.port": "5432",
    "database.user": "asyura",
    "database.password": "$Mahesvara14",
    "database.dbname": "ecommerce_product",
    "topic.prefix": "ecommerce",
    "plugin.name": "pgoutput",
    "table.include.list": "public.products",
    "slot.name": "product_slot",
    "publication.autocreate.mode": "filtered"
  }
}'

# Mengecek status konektor, pastikan statusnya menampilkan "RUNNING"
curl -s http://localhost:8083/connectors/product-connector/status

# Menjalankan server aplikasi Golang untuk menerima HTTP request
cd product-service
make run

# Run Search Service Worker (Event Consumer)
cd search-service
go run cmd/worker/main.go

# Run Search Service API (Query Side)
cd search-service
go run cmd/api/main.go