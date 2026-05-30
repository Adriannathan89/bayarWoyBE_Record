package user_profile_service

import (
	"bayar-woy-project/config"
	"bayar-woy-project/models"
	"bayar-woy-project/responses"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetProfile(c *gin.Context) {
	userID := c.GetString("userID")

	var user models.User
	if err := config.DB.First(&user, "id = ?", userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "User not found"})
		return
	}

	discordInfo := responses.DiscordProfileInfo{
		Connected:          user.IsValidated && user.DiscordID != nil,
		CommitNotifEnabled: user.DiscordCommitNotifEnabled,
		WeeklyNotifEnabled: user.DiscordWeeklyNotifEnabled,
	}
	if user.DiscordUsername != nil {
		discordInfo.Username = *user.DiscordUsername
	}

	c.JSON(http.StatusOK, responses.APIResponse{
		StatusCode: http.StatusOK,
		Message:    "OK",
		Data: responses.ProfileResponse{
			ID:       user.ID,
			Username: user.Username,
			Discord:  discordInfo,
		},
	})
}
