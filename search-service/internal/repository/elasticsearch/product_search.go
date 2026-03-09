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

func (r *productSearchRepo) SearchProducts(ctx context.Context, query string) ([]domain.Product, error) {
	// Build the Elasticsearch Match Query
	esQuery := fmt.Sprintf(`{
		"query": {
			"match": {
				"name": "%s"
			}
		}
	}`, query)

	// Execute the search
	res, err := r.esClient.Search(
		r.esClient.Search.WithContext(ctx),
		r.esClient.Search.WithIndex("products"),
		r.esClient.Search.WithBody(strings.NewReader(esQuery)),
	)
	if err != nil {
		return nil, fmt.Errorf("[ERROR] Failed to execute search: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		return nil, fmt.Errorf("[ERROR] Error response from search engine: %s", res.Status())
	}

	// Decode the response
	var raw map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&raw); err != nil {
		return nil, fmt.Errorf("[ERROR] Failed to parse search results: %w", err)
	}

	// Extract the products
	var products []domain.Product
	hits := raw["hits"].(map[string]interface{})["hits"].([]interface{})

	for _, hit := range hits {
		source := hit.(map[string]interface{})["_source"]
		sourceBytes, _ := json.Marshal(source)

		var prod domain.Product
		json.Unmarshal(sourceBytes, &prod)
		products = append(products, prod)
	}

	return products, nil
}
