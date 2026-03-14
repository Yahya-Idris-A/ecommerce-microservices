package elasticsearch

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/Yahya-idris-A/ecommerce-microservices/search-service/internal/domain"
	"github.com/elastic/go-elasticsearch/v8"
)

type productSearchRepo struct {
	esClient *elasticsearch.Client
}

func NewProductSearchRepository(es *elasticsearch.Client) domain.ProductSearchRepository {
	return &productSearchRepo{
		esClient: es,
	}
}

func (r *productSearchRepo) SearchGlobal(ctx context.Context, keyword string, limit int, cursor string) (*domain.GlobalSearchResult, string, error) {
	// 1. Susun blok pencarian dasar dan urutannya (Tie-breaker ID wajib ada)
	esQuery := fmt.Sprintf(`{
        "size": %d,
        "query": {
            "multi_match": {
                "query": "%s",
                "fields": ["name", "description"]
            }
        },
        "sort": [
            { "_score": "desc" },
            { "id": "asc" }
        ]`, limit, keyword)

	if cursor != "" {
		parts := strings.Split(cursor, "|")
		if len(parts) == 2 {
			// parts[0] adalah skor (angka), parts[1] adalah _id (string)
			esQuery += fmt.Sprintf(`, "search_after": [%s, "%s"]`, parts[0], parts[1])
		}
	}
	esQuery += `}`

	res, err := r.esClient.Search(
		r.esClient.Search.WithContext(ctx),
		r.esClient.Search.WithIndex("products", "merchants"),
		r.esClient.Search.WithBody(strings.NewReader(esQuery)),
	)
	if err != nil {
		return nil, "", err
	}
	defer res.Body.Close()

	if res.IsError() {
		var e map[string]interface{}
		if err := json.NewDecoder(res.Body).Decode(&e); err != nil {
			return nil, "", fmt.Errorf("[ES ERROR] Gagal membaca pesan error dari ES: %s", err)
		}
		// Ini akan melempar pesan error asli dari Elasticsearch ke Postman
		return nil, "", fmt.Errorf("[ES ERROR] %v", e["error"])
	}

	var raw map[string]interface{}
	json.NewDecoder(res.Body).Decode(&raw)

	result := &domain.GlobalSearchResult{
		Products:  []domain.Product{},
		Merchants: []domain.Merchant{},
	}

	hitsMap, ok := raw["hits"].(map[string]interface{})
	if !ok {
		return result, "", nil
	}
	hits := hitsMap["hits"].([]interface{})

	for _, hit := range hits {
		hitDict := hit.(map[string]interface{})
		index := hitDict["_index"].(string)
		sourceBytes, _ := json.Marshal(hitDict["_source"])

		if index == "products" {
			var prod domain.Product
			json.Unmarshal(sourceBytes, &prod)
			result.Products = append(result.Products, prod)
		} else if index == "merchants" {
			var merch domain.Merchant
			json.Unmarshal(sourceBytes, &merch)
			result.Merchants = append(result.Merchants, merch)
		}
	}

	var nextCursor string
	if len(hits) > 0 {
		lastHit := hits[len(hits)-1].(map[string]interface{})
		if sortValues, hasSort := lastHit["sort"].([]interface{}); hasSort && len(sortValues) == 2 {
			// Gabungkan score dan _id menjadi format "score|_id"
			nextCursor = fmt.Sprintf("%v|%v", sortValues[0], sortValues[1])
		}
	}

	return result, nextCursor, nil
}
