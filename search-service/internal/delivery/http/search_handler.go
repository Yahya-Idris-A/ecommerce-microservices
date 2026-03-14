package http

import (
	"net/http"
	"strconv"

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
	keyword := c.Query("q")
	cursor := c.Query("cursor")

	// Validasi Limit Anti-Bug
	limitStr := c.Query("limit")
	limit := 10 // Default
	if limitStr != "" {
		if parsedLimit, err := strconv.Atoi(limitStr); err == nil && parsedLimit > 0 {
			limit = parsedLimit
		}
	}

	// Call the Usecase layer
	result, nextCursor, err := h.searchUsecase.Search(c.Request.Context(), keyword, limit, cursor)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Sisipkan metadata persis seperti di Product Service
	c.JSON(http.StatusOK, gin.H{
		"message": "Global search successful",
		"data":    result,
		"meta": gin.H{
			"next_cursor": nextCursor,
			"has_more":    nextCursor != "",
		},
	})
}
