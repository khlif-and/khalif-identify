package middleware

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/redis/go-redis/v9"

)

// Update: Menambahkan parameter Redis Client (rdb)
func AuthMiddleware(secretKey string, rdb *redis.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if !strings.Contains(authHeader, "Bearer") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}

		tokenString := strings.Split(authHeader, " ")[1]

		// --- LOGIC BLACKLIST CHECK ---
		// Cek apakah token ini ada di daftar blacklist Redis?
		ctx := context.Background()
		_, err := rdb.Get(ctx, "blacklist:"+tokenString).Result()
		if err == nil {
			// Jika ditemukan di Redis (err == nil), berarti token sudah logout/hangus
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Token invalidated (Logged out)"})
			return
		}
		// -----------------------------

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method")
			}
			return []byte(secretKey), nil
		})

		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid Token"})
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			c.Set("user_id", claims["user_id"])
			c.Set("role", claims["role"])
			c.Next()
		} else {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid Token Claims"})
		}
	}
}