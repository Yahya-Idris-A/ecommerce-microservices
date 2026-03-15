package http

import (
	"context"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	// Ingat sesuaikan path import ini
	"github.com/Yahya-idris-A/ecommerce-microservices/product-service/internal/delivery/http/middleware"
	"github.com/Yahya-idris-A/ecommerce-microservices/product-service/internal/domain"
)

// ProductHandler bertugas menjembatani HTTP request dengan Usecase
type ProductHandler struct {
	usecase domain.ProductUsecase
}

// NewProductHandler adalah constructor untuk inisialisasi routing.
// Kita menerima engine Gin dan menyuntikkan Usecase ke dalam handler.
func NewProductHandler(r *gin.Engine, us domain.ProductUsecase) {
	handler := &ProductHandler{
		usecase: us,
	}

	// Membuat grup routing untuk versi API
	api := r.Group("/api/v1/products")
	{
		// Mendaftarkan endpoint ke fungsi yang sesuai
		api.POST("", middleware.AuthGuard("admin", "merchant"), handler.Create)
		api.GET("", handler.GetProducts)
		api.GET("/me", middleware.AuthGuard("merchant"), handler.GetMyProducts)
		api.GET("/:id", handler.GetByID)
		api.DELETE("/:id", middleware.AuthGuard("merchant"), handler.DeleteProduct)
	}
}

// Create menangani request POST untuk membuat produk baru
func (h *ProductHandler) Create(c *gin.Context) {
	var req domain.CreateProductRequest

	// 1. Parsing JSON dari request body ke struct DTO
	if err := c.ShouldBindJSON(&req); err != nil {
		// Mengembalikan 400 Bad Request jika format JSON salah atau data tidak lengkap
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid request payload",
			"details": err.Error(),
		})
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()
	// 2. Memanggil Usecase (Koki) untuk memproses data
	product, err := h.usecase.CreateProduct(ctx, &req)
	if err != nil {
		log.Printf("[DEBUG DB ERROR]: %v\n", err)
		// Mengembalikan 500 Internal Server Error jika terjadi kegagalan di layer bisnis/database
		// Catatan: Di sistem aslinya, kita harus memisahkan error bisnis (4xx) dan error server (5xx).
		// Kita akan buat sistem Custom Error yang rapi nanti.
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to process request",
		})
		return
	}

	// 3. Mengembalikan respons sukses (201 Created)
	c.JSON(http.StatusCreated, gin.H{
		"message": "product created successfully",
		"data":    product,
	})
}

// GetByID menangani request GET untuk mencari produk
func (h *ProductHandler) GetByID(c *gin.Context) {
	// Mengambil parameter 'id' dari URL (contoh: /api/v1/products/123-abc)
	idParam := c.Param("id")

	// Validasi apakah format string adalah UUID yang valid
	productID, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid product id format",
		})
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()

	// Meminta Usecase untuk mencari data
	product, err := h.usecase.GetProductByID(ctx, productID)
	if err != nil {
		// Asumsi error dari Usecase saat ini adalah "data tidak ditemukan"
		c.JSON(http.StatusNotFound, gin.H{
			"error": "product not found",
		})
		return
	}

	// Mengembalikan respons sukses (200 OK)
	c.JSON(http.StatusOK, gin.H{
		"message": "success",
		"data":    product,
	})
}

// Tambahkan fungsi ini di bagian bawah file product_handler.go
func (h *ProductHandler) GetProducts(c *gin.Context) {
	merchantID := c.Query("merchant_id")
	keyword := c.Query("q")
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

	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()

	products, nextCursor, err := h.usecase.GetAllProducts(ctx, merchantID, keyword, limit, cursorCreatedAt, cursorID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch products"})
		return
	}

	// Jika produk kosong, kembalikan array kosong agar frontend Next.js kamu tidak error saat map()
	if products == nil {
		products = []domain.Product{}
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Product retrieved successfully",
		"data":    products,
		"meta": gin.H{
			"next_cursor": nextCursor,
			"has_more":    nextCursor != "",
		},
	})
}

func (h *ProductHandler) GetMyProducts(c *gin.Context) {
	// Ambil user_id dari JWT token yang sudah dicegat oleh AuthGuard
	userIDStr, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	userID, err := uuid.Parse(userIDStr.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid User ID format"})
		return
	}

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
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid cursor format"})
			return
		}
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()

	// Tarik data lewat Usecase
	products, nextCursor, err := h.usecase.GetMyProducts(ctx, userID, limit, cursorCreatedAt, cursorID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	if products == nil {
		products = []domain.Product{}
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "My catalog retrieved successfully",
		"data":    products,
		"meta": gin.H{
			"next_cursor": nextCursor,
			"has_more":    nextCursor != "",
		},
	})
}

func (h *ProductHandler) DeleteProduct(c *gin.Context) {
	userIDStr, _ := c.Get("user_id")
	userID, _ := uuid.Parse(userIDStr.(string))

	idParam := c.Param("id")

	// Validasi apakah format string adalah UUID yang valid
	productID, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid product id format",
		})
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()

	err = h.usecase.DeleteProduct(ctx, userID, productID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Product deleted successfully",
	})
}
