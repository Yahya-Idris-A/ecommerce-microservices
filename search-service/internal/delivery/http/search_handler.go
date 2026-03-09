package http

import (
	"net/http"

	"github.com/Yahya-idris-A/ecommerce-microservices/search-service/internal/domain"
	"github.com/gin-gonic/gin"
)

type SearchHandler struct {
	searchUsecase domain.ProductSearchUsecase
}

func NewSearchHandler(r *gin.Engine, us domain.ProductSearchUsecase) {
	handler := &SearchHandler{
		searchUsecase: us,
	}

	// Define the routing group
	api := r.Group("/api/v1")
	{
		api.GET("/search", handler.SearchProducts)
	}
}

func (h *SearchHandler) SearchProducts(c *gin.Context) {
	query := c.Query("q")

	// Call the Usecase layer
	products, err := h.searchUsecase.Search(c.Request.Context(), query)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Return clean response
	c.JSON(http.StatusOK, gin.H{
		"message": "Search successful",
		"count":   len(products),
		"data":    products,
	})
}
