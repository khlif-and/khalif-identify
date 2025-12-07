package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"khalif-identify/internal/domain"

)

// MaxAdminLimit mencegah registrasi jika jumlah admin sudah mencapai batas
func MaxAdminLimit(repo domain.UserRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		const MaxAdminCount = 3
		const AdminRoleID = 1

		// Perbaikan: Tambahkan Context dari request
		currentCount, err := repo.CountByRoleID(c.Request.Context(), AdminRoleID)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to check admin quota"})
			return
		}

		if currentCount >= MaxAdminCount {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"error": "Admin quota is full (Max 3 admins allowed)",
			})
			return
		}

		c.Next()
	}
}