package guard

import (
	"bayar-woy-project/config"
	"bayar-woy-project/models"
	"time"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString, err := c.Cookie("token")

		claims, err := ValidateAccessToken(tokenString)
		if err == nil {
			c.Set("userID", claims.UserID)
			c.Set("username", claims.Username)
			c.Next()
			return
		}
		refreshToken, err := c.Cookie("refresh_token")
		if err != nil {
			c.JSON(401, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}

		refreshClaims, err := ValidateRefreshToken(refreshToken)
		if err != nil {
			c.JSON(401, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}
		var sesion models.Session

		if err := config.DB.Where("refresh_token = ?", refreshToken).First(&sesion).Error; err != nil {
			c.JSON(401, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}

		if sesion.ExpiresAt.Before(time.Now()) {
			c.JSON(401, gin.H{"error": "Session expired"})
			c.Abort()
			return
		}

		newAccessToken, _ := GenerateToken(refreshClaims.Username, refreshClaims.UserID)
		c.SetCookie("token", newAccessToken, 60*10, "/", config.GetEnv("CLIENT_HOST"), true, true)

		c.Set("userID", refreshClaims.UserID)
		c.Set("username", refreshClaims.Username)
		c.Next()
	}
}
