package user_profile_service

import (
	"bayar-woy-project/config"
	"bayar-woy-project/models"
	"bayar-woy-project/responses"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetDiscordStatus(c *gin.Context) {
	userID := c.GetString("userID")

	var user models.User
	if err := config.DB.First(&user, "id = ?", userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "User not found"})
		return
	}

	resp := responses.DiscordStatusResponse{
		Verified: user.IsValidated && user.DiscordID != nil,
	}
	if user.DiscordUsername != nil {
		resp.DiscordUsername = *user.DiscordUsername
	}

	c.JSON(http.StatusOK, responses.APIResponse{
		StatusCode: http.StatusOK,
		Message:    "OK",
		Data:       resp,
	})
}
