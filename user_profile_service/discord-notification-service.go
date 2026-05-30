package user_profile_service

import (
	"bayar-woy-project/config"
	"bayar-woy-project/dto"
	"bayar-woy-project/models"
	"bayar-woy-project/responses"
	"net/http"

	"github.com/gin-gonic/gin"
)

func UpdateNotifSettings(c *gin.Context) {
	userID := c.GetString("userID")

	var req dto.ToggleNotificationDto
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Bad Request"})
		return
	}

	column := ""
	switch req.Type {
	case "commit":
		column = "discord_commit_notif_enabled"
	case "weekly":
		column = "discord_weekly_notif_enabled"
	default:
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid type"})
		return
	}

	if err := config.DB.Model(&models.User{}).Where("id = ?", userID).
		Update(column, req.Enabled).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to update"})
		return
	}

	c.JSON(http.StatusOK, responses.APIResponse{
		StatusCode: http.StatusOK,
		Message:    "Notification settings updated",
		Data:       gin.H{"success": true},
	})
}
