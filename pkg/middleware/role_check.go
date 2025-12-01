package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"

)

func OnlyAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get("role")
		if !exists {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}

		if role != "Admin" {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"error": "Access denied. Admins only.",
			})
			return
		}

		c.Next()
	}
}