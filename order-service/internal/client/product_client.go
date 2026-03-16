package client

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/Yahya-idris-A/ecommerce-microservices/order-service/internal/usecase"
	"github.com/google/uuid"
)

type productClient struct {
	httpClient *http.Client
	productURL string // Contoh: http://localhost:8081 (Port Product Service)
	userURL    string // Contoh: http://localhost:8080 (Port User Service)
}

// Pastikan tipe yang dikembalikan adalah interface dari Usecase
func NewProductClient(productURL string, userURL string) usecase.ProductServiceClient {
	return &productClient{
		// Kita pasang timeout global 5 detik untuk client HTTP ini
		// Biar kalau product-service mati, order-service nggak ikut hang lama-lama
		httpClient: &http.Client{Timeout: 5 * time.Second},
		productURL: productURL,
		userURL:    userURL,
	}
}

// 1. Mengambil Detail Produk (Nembak Product Service)
func (c *productClient) GetProductDetails(ctx context.Context, productID uuid.UUID) (string, string, float64, error) {
	url := fmt.Sprintf("%s/api/v1/products/%s", c.productURL, productID.String())

	// Buat request dengan membawa context (Stopwatch kita!)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return "", "", 0, err
	}

	// Eksekusi request
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", "", 0, fmt.Errorf("failed to fetch product: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", "", 0, errors.New("product not found or error in product service")
	}

	// Parsing JSON dari Product Service
	// Asumsi response JSON Product Service-mu punya format:
	// {"data": {"name": "Laptop", "image_url": "...", "price": 15000000}}
	var result struct {
		Data struct {
			Name     string  `json:"name"`
			ImageURL string  `json:"image_url"`
			Price    float64 `json:"price"`
		} `json:"data"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", "", 0, err
	}

	return result.Data.Name, result.Data.ImageURL, result.Data.Price, nil
}

// 2. Mengambil Nama Toko/Merchant (Nembak User Service)
func (c *productClient) GetMerchantName(ctx context.Context, merchantID uuid.UUID) (string, error) {
	// Karena merchant sebenarnya adalah entitas User dengan role 'merchant',
	// kita nembak ke User Service untuk ambil nama lengkapnya.
	url := fmt.Sprintf("%s/api/v1/users/profile/%s", c.userURL, merchantID.String())

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return "", err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to fetch merchant: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		// Kalau gagal, kita return default string saja biar keranjang nggak error
		return "Unknown Merchant", nil
	}

	var result struct {
		Data struct {
			FullName string `json:"full_name"`
		} `json:"data"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "Unknown Merchant", nil
	}

	return result.Data.FullName, nil
}
