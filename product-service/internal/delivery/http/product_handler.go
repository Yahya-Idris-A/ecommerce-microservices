package http

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	// Ingat sesuaikan path import ini
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
		api.POST("", handler.Create)
		api.GET("/:id", handler.GetByID)
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

	// 2. Memanggil Usecase (Koki) untuk memproses data
	product, err := h.usecase.CreateProduct(&req)
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

	// Meminta Usecase untuk mencari data
	product, err := h.usecase.GetProductByID(productID)
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
