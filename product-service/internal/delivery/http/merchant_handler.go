package http

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/Yahya-idris-A/ecommerce-microservices/product-service/internal/delivery/http/middleware"
	"github.com/Yahya-idris-A/ecommerce-microservices/product-service/internal/domain"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type MerchantHandler struct {
	merchantUsecase domain.MerchantUsecase
}

func NewMerchantHandler(r *gin.Engine, us domain.MerchantUsecase) {
	handler := &MerchantHandler{
		merchantUsecase: us,
	}

	api := r.Group("/api/v1")
	{
		api.GET("/merchants", handler.GetMerchants)
		api.GET("/merchants/:id", handler.GetMerchantByID)
		api.GET("/merchants/me", middleware.AuthGuard("merchant"), handler.GetMyProfile)
		// Hanya user dengan role "merchant" yang boleh membuat profil toko
		api.POST("/merchants", middleware.AuthGuard("merchant"), handler.CreateMerchant)
	}
}

func (h *MerchantHandler) CreateMerchant(c *gin.Context) {
	// 1. Ambil ID User dari Token JWT (disuntikkan oleh AuthGuard)
	userIDStr, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found in token"})
		return
	}

	// Parsing string UUID ke tipe uuid.UUID
	userID, err := uuid.Parse(userIDStr.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid User ID format in token"})
		return
	}

	// 2. Bind payload JSON
	var req domain.CreateMerchantRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	// 3. Eksekusi Usecase
	merchant, err := h.merchantUsecase.CreateMerchant(userID, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create merchant profile"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Merchant profile created successfully",
		"data":    merchant,
	})
}

func (h *MerchantHandler) GetMyProfile(c *gin.Context) {
	// 1. Ambil user_id dari JWT token (yang di-set oleh middleware AuthGuard)
	userIDStr, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized: No user ID in token"})
		return
	}

	// 2. Parse string menjadi UUID
	userID, err := uuid.Parse(userIDStr.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid User ID format in token"})
		return
	}

	// 3. Tarik data profil toko lewat Usecase
	merchant, err := h.merchantUsecase.GetByUserID(userID)
	if err != nil {
		// Jika errornya adalah "merchant profile not found" dari Repo
		if err.Error() == "merchant profile not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": "Anda belum membuat toko"})
			return
		}
		// Jika error sistem lainnya
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 4. Kembalikan data toko
	c.JSON(http.StatusOK, gin.H{
		"message": "Merchant profile retrieved successfully",
		"data":    merchant,
	})
}

// GetMerchants menangani Get All sekaligus Search (jika ada parameter ?q=)
func (h *MerchantHandler) GetMerchants(c *gin.Context) {
	keyword := c.Query("q") // Tangkap parameter dari URL
	limitStr := c.Query("limit")
	limit := 10 // Default limit jika user tidak menyertakan parameter ini

	if limitStr != "" {
		parsedLimit, err := strconv.Atoi(limitStr)
		// Hanya gunakan nilai dari user jika dia angka valid dan lebih dari 0
		if err == nil && parsedLimit > 0 {
			limit = parsedLimit
		}
	}
	cursor := c.Query("cursor")

	var cursorCreatedAt, cursorID string

	if cursor != "" {
		parts := strings.Split(cursor, "|")
		if len(parts) == 2 {
			cursorCreatedAt = parts[0]
			cursorID = parts[1]
		} else {
			// Jika formatnya bukan "waktu|id", tolak requestnya
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid cursor format"})
			return
		}
	}

	merchants, nextCursor, err := h.merchantUsecase.GetAllMerchants(keyword, limit, cursorCreatedAt, cursorID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch merchants"})
		return
	}

	if merchants == nil {
		merchants = []domain.Merchant{}
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Merchants retrieved successfully",
		"data":    merchants,
		"meta": gin.H{
			"next_cursor": nextCursor,
			"has_more":    nextCursor != "",
		},
	})
}

// GetMerchantByID menangani pengambilan profil toko tunggal
func (h *MerchantHandler) GetMerchantByID(c *gin.Context) {
	idParam := c.Param("id") // Tangkap ID dari path URL (contoh: /merchants/123-abc)

	id, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid merchant ID format"})
		return
	}

	merchant, err := h.merchantUsecase.GetMerchantByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Merchant retrieved successfully",
		"data":    merchant,
	})
}
