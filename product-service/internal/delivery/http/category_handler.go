package http

import (
	"net/http"

	"github.com/google/uuid"

	"github.com/Yahya-idris-A/ecommerce-microservices/product-service/internal/delivery/http/middleware"
	"github.com/Yahya-idris-A/ecommerce-microservices/product-service/internal/domain"
	"github.com/gin-gonic/gin"
)

type CategoryHandler struct {
	categoryUsecase domain.CategoryUsecase
	merchantRepo    domain.MerchantRepository
}

func NewCategoryHandler(r *gin.Engine, us domain.CategoryUsecase, mr domain.MerchantRepository) {
	handler := &CategoryHandler{
		categoryUsecase: us,
		merchantRepo:    mr,
	}

	api := r.Group("/api/v1")
	{
		// Kita lindungi dengan AuthGuard
		api.POST("/categories", middleware.AuthGuard("admin", "merchant"), handler.CreateCategory)
		api.GET("/categories", handler.GetAllCategories)
		api.GET("/categories/:id", handler.GetByID)
	}
}

func (h *CategoryHandler) CreateCategory(c *gin.Context) {
	var req domain.CreateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	userRole, _ := c.Get("user_role")
	userIDStr, _ := c.Get("user_id")

	var merchantIDPtr *uuid.UUID

	if userRole == "merchant" {
		userID, err := uuid.Parse(userIDStr.(string))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid User ID"})
			return
		}

		merchant, err := h.merchantRepo.GetByUserID(userID)
		if err != nil {
			// Tolak jika dia merchant tapi belum buka profil toko
			c.JSON(http.StatusForbidden, gin.H{"error": "Merchant profile not found. Please create a store first."})
			return
		}

		// Set pointer ke ID toko yang ditemukan
		merchantIDPtr = &merchant.ID
	}

	category, err := h.categoryUsecase.CreateCategory(&req, merchantIDPtr)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create category"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Category created successfully",
		"data":    category,
	})
}

// GetByID menangani request GET untuk mencari kategori
func (h *CategoryHandler) GetByID(c *gin.Context) {
	// Mengambil parameter 'id' dari URL (contoh: /api/v1/categories/123-abc)
	idParam := c.Param("id")

	// Validasi apakah format string adalah UUID yang valid
	categoryID, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid category id format",
		})
		return
	}

	// Meminta Usecase untuk mencari data
	category, err := h.categoryUsecase.GetCategoryByID(categoryID)
	if err != nil {
		// Asumsi error dari Usecase saat ini adalah "data tidak ditemukan"
		c.JSON(http.StatusNotFound, gin.H{
			"error": "category not found",
		})
		return
	}

	// Mengembalikan respons sukses (200 OK)
	c.JSON(http.StatusOK, gin.H{
		"message": "success",
		"data":    category,
	})
}

func (h *CategoryHandler) GetAllCategories(c *gin.Context) {
	merchantID := c.Query("merchant_id")
	categories, err := h.categoryUsecase.GetAllCategories(merchantID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch categories"})
		return
	}

	// Jika kategori kosong, kembalikan array kosong agar frontend Next.js kamu tidak error saat map()
	if categories == nil {
		categories = []domain.Category{}
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Categories retrieved successfully",
		"data":    categories,
	})
}
