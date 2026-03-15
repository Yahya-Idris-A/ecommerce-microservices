package http

import (
	"context"
	"net/http"
	"time"

	"github.com/Yahya-idris-A/ecommerce-microservices/user-service/internal/delivery/http/middleware"
	"github.com/Yahya-idris-A/ecommerce-microservices/user-service/internal/domain"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type AddressHandler struct {
	addressUsecase domain.AddressUsecase
}

func NewAddressHandler(r *gin.Engine, addressUsecase domain.AddressUsecase) {
	addressHandler := &AddressHandler{
		addressUsecase: addressUsecase,
	}

	// Membuat grup routing untuk versi API
	protectedUserRoutes := r.Group("/api/v1/users/me")
	protectedUserRoutes.Use(middleware.AuthGuard("admin", "merchant", "buyer"))
	{

		// Manajemen Buku Alamat
		protectedUserRoutes.POST("/addresses", addressHandler.AddAddress)
		protectedUserRoutes.GET("/addresses", addressHandler.GetMyAddresses)
		protectedUserRoutes.PUT("/addresses/:id", addressHandler.UpdateAddress)
		protectedUserRoutes.DELETE("/addresses/:id", addressHandler.DeleteAddress)
		protectedUserRoutes.PUT("/addresses/:id/primary", addressHandler.SetPrimaryAddress)
	}
}

func (h *AddressHandler) AddAddress(c *gin.Context) {
	userIDStr, _ := c.Get("user_id")
	userID, _ := uuid.Parse(userIDStr.(string))

	var req domain.CreateAddressReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()

	address, err := h.addressUsecase.AddAddress(ctx, userID, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Address added successfully",
		"data":    address,
	})
}

func (h *AddressHandler) GetMyAddresses(c *gin.Context) {
	userIDStr, _ := c.Get("user_id")
	userID, _ := uuid.Parse(userIDStr.(string))

	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()

	addresses, err := h.addressUsecase.GetMyAddresses(ctx, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Pastikan return array kosong jika tidak ada alamat, bukan null
	if addresses == nil {
		addresses = []domain.Address{}
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Addresses retrieved successfully",
		"data":    addresses,
	})
}

func (h *AddressHandler) UpdateAddress(c *gin.Context) {
	userIDStr, _ := c.Get("user_id")
	userID, _ := uuid.Parse(userIDStr.(string))

	addressIDStr := c.Param("id")
	addressID, err := uuid.Parse(addressIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid address ID"})
		return
	}

	var req domain.UpdateAddressReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()

	address, err := h.addressUsecase.UpdateAddress(ctx, userID, addressID, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Address updated successfully",
		"data":    address,
	})
}

func (h *AddressHandler) DeleteAddress(c *gin.Context) {
	userIDStr, _ := c.Get("user_id")
	userID, _ := uuid.Parse(userIDStr.(string))

	addressIDStr := c.Param("id")
	addressID, err := uuid.Parse(addressIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid address ID"})
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()

	err = h.addressUsecase.DeleteAddress(ctx, userID, addressID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Address deleted successfully",
	})
}

func (h *AddressHandler) SetPrimaryAddress(c *gin.Context) {
	userIDStr, _ := c.Get("user_id")
	userID, _ := uuid.Parse(userIDStr.(string))

	addressIDStr := c.Param("id")
	addressID, err := uuid.Parse(addressIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid address ID"})
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()

	err = h.addressUsecase.SetPrimaryAddress(ctx, userID, addressID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Primary address updated successfully",
	})
}
