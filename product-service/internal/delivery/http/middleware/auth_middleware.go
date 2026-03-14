package middleware

import (
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// AuthGuard protects routes and checks if the user has the required roles
func AuthGuard(allowedRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. Get the Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "[ERROR] Authorization header is required"})
			return
		}

		// 2. Check if the format is "Bearer <token>"
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "[ERROR] Invalid authorization header format"})
			return
		}

		tokenString := parts[1]
		secret := os.Getenv("JWT_SECRET")

		// 3. Parse and validate the token signature
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Ensure the signing method is HMAC (HS256)
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, http.ErrAbortHandler
			}
			return []byte(secret), nil
		})

		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "[ERROR] Invalid or expired token"})
			return
		}

		// 4. Extract claims (payload)
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "[ERROR] Invalid token payload"})
			return
		}

		userRole := claims["role"].(string)
		userID := claims["user_id"].(string)

		// 5. Check if the user's role is in the allowedRoles list
		roleAllowed := false
		for _, role := range allowedRoles {
			if userRole == role {
				roleAllowed = true
				break
			}
		}

		if !roleAllowed {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "[ERROR] You do not have permission to access this resource"})
			return
		}

		// 6. If valid, save the user info to the context so the handler can use it later
		c.Set("user_id", userID)
		c.Set("user_role", userRole)

		// Continue to the actual handler (e.g., Create Product)
		c.Next()
	}
}
