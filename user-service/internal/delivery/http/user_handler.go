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

type UserHandler struct {
	userUsecase domain.UserUsecase
}

func NewUserHandler(r *gin.Engine, userUsecase domain.UserUsecase) {
	handler := &UserHandler{
		userUsecase: userUsecase,
	}

	// Membuat grup routing untuk versi API
	api := r.Group("/api/v1/users")
	{
		api.POST("/register", handler.Register)
		api.POST("/login", handler.Login)
	}
	protectedUserRoutes := r.Group("/api/v1/users/me")
	protectedUserRoutes.Use(middleware.AuthGuard("admin", "merchant", "buyer"))
	{
		// Manajemen Profil
		protectedUserRoutes.GET("", handler.GetProfile)
		protectedUserRoutes.PUT("", handler.UpdateProfile)
		protectedUserRoutes.DELETE("", handler.DeleteAccount)
	}
}

func (h *UserHandler) Register(c *gin.Context) {
	var req domain.RegisterReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()

	user, err := h.userUsecase.Register(ctx, &req)
	if err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "User registered successfully",
		"data":    user,
	})
}

func (h *UserHandler) Login(c *gin.Context) {
	var req domain.LoginReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()

	res, err := h.userUsecase.Login(ctx, &req)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Login successful",
		"data":    res,
	})
}

func (h *UserHandler) GetProfile(c *gin.Context) {
	userIDStr, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	userID, _ := uuid.Parse(userIDStr.(string))

	// 1. BIKIN STOPWATCH-NYA DI SINI BRO!
	// Kita ambil context bawaan dari request Gin, lalu pasang timer 10 detik
	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)

	// 2. WAJIB PANGGIL DEFER CANCEL
	// Ini untuk membersihkan memori stopwatch-nya saat fungsi GetProfile selesai,
	// entah itu selesai karena sukses, atau selesai karena kena error.
	defer cancel()

	user, err := h.userUsecase.GetProfile(ctx, userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Profile retrieved successfully",
		"data":    user,
	})
}

func (h *UserHandler) UpdateProfile(c *gin.Context) {
	userIDStr, _ := c.Get("user_id")
	userID, _ := uuid.Parse(userIDStr.(string))

	var req domain.UpdateProfileReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()

	user, err := h.userUsecase.UpdateProfile(ctx, userID, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Profile updated successfully",
		"data":    user,
	})
}

func (h *UserHandler) DeleteAccount(c *gin.Context) {
	userIDStr, _ := c.Get("user_id")
	userID, _ := uuid.Parse(userIDStr.(string))

	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()

	err := h.userUsecase.DeleteAccount(ctx, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Account deleted successfully",
	})
}
