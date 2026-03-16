package http

import (
	"context"
	"net/http"
	"time"

	"github.com/Yahya-idris-A/ecommerce-microservices/order-service/internal/delivery/http/middleware"
	"github.com/Yahya-idris-A/ecommerce-microservices/order-service/internal/domain"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type CartHandler struct {
	cartUsecase domain.CartUsecase
}

// Perhatikan kita langsung mendaftarkan route di dalam constructor
func NewCartHandler(r *gin.Engine, cartUsecase domain.CartUsecase, authMiddleware gin.HandlerFunc) {
	handler := &CartHandler{cartUsecase: cartUsecase}

	// Semua endpoint keranjang butuh token (Harus Login)
	cartRoutes := r.Group("/api/v1/carts")
	cartRoutes.Use(middleware.AuthGuard("admin", "merchant", "buyer"))
	{
		cartRoutes.POST("/items", handler.AddItem)
		cartRoutes.GET("", handler.GetMyCart)
		cartRoutes.PUT("/items/:id", handler.UpdateItemQuantity)
		cartRoutes.DELETE("/items/:id", handler.RemoveItemFromCart)
	}
}

func (h *CartHandler) AddItem(c *gin.Context) {
	// 1. Ambil User ID dari Token (via Middleware)
	userIDStr, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	userID, _ := uuid.Parse(userIDStr.(string))

	// 2. Validasi Body Request
	var req domain.AddCartItemReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 3. Pasang Stopwatch 5 Detik
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	// 4. Lempar ke Usecase
	err := h.cartUsecase.AddItemToCart(ctx, userID, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Item added to cart successfully",
	})
}

func (h *CartHandler) GetMyCart(c *gin.Context) {
	userIDStr, _ := c.Get("user_id")
	userID, _ := uuid.Parse(userIDStr.(string))

	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	cart, err := h.cartUsecase.GetMyCart(ctx, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Cart retrieved successfully",
		"data":    cart,
	})
}

func (h *CartHandler) UpdateItemQuantity(c *gin.Context) {
	userIDStr, _ := c.Get("user_id")
	userID, _ := uuid.Parse(userIDStr.(string))

	itemIDStr := c.Param("id")
	itemID, err := uuid.Parse(itemIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid item id format"})
		return
	}

	var req domain.UpdateCartItemReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	err = h.cartUsecase.UpdateItemQuantity(ctx, userID, itemID, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Cart item updated successfully",
	})
}

func (h *CartHandler) RemoveItemFromCart(c *gin.Context) {
	userIDStr, _ := c.Get("user_id")
	userID, _ := uuid.Parse(userIDStr.(string))

	itemIDStr := c.Param("id")
	itemID, err := uuid.Parse(itemIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid item id format"})
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	err = h.cartUsecase.RemoveItemFromCart(ctx, userID, itemID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Item removed from cart successfully",
	})
}
