package user_profile_service

import (
	"bayar-woy-project/config"
	"bayar-woy-project/models"
	"bayar-woy-project/responses"
	"net/http"

	"github.com/gin-gonic/gin"
)

func DisconnectDiscord(c *gin.Context) {
	userID := c.GetString("userID")

	if err := config.DB.Model(&models.User{}).Where("id = ?", userID).
		Updates(map[string]interface{}{
			"discord_id":       nil,
			"discord_username": nil,
			"is_validated":     false,
		}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to disconnect"})
		return
	}

	config.DB.Where("user_id = ?", userID).Delete(&models.DiscordVerification{})

	c.JSON(http.StatusOK, responses.APIResponse{
		StatusCode: http.StatusOK,
		Message:    "Discord disconnected",
		Data:       gin.H{"success": true},
	})
}
