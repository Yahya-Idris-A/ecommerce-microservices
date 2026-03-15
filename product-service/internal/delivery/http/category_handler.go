package http

import (
	"context"
	"net/http"
	"time"

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

	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()

	category, err := h.categoryUsecase.CreateCategory(ctx, &req, merchantIDPtr)
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

	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()

	// Meminta Usecase untuk mencari data
	category, err := h.categoryUsecase.GetCategoryByID(ctx, categoryID)
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

	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()

	categories, err := h.categoryUsecase.GetAllCategories(ctx, merchantID)
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

func (h *CategoryHandler) DeleteCategory(c *gin.Context) {
	roleStr, exists := c.Get("user_role")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized: No role found in token"})
		return
	}
	role := roleStr.(string)

	userIDStr, _ := c.Get("user_id")

	var merchantIDPtr *uuid.UUID

	if role == "merchant" {
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

	categoryIDStr := c.Param("id")
	categoryID, err := uuid.Parse(categoryIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid category id format"})
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()

	err = h.categoryUsecase.DeleteCategory(ctx, role, *merchantIDPtr, categoryID)
	if err != nil {
		// Cek apakah errornya karena masalah otorisasi
		if err.Error() == "unauthorized: only admin can delete global categories" ||
			err.Error() == "unauthorized: category does not belong to this merchant" {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Category deleted successfully",
	})
}
