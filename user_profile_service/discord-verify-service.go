package user_profile_service

import (
	"bayar-woy-project/config"
	"bayar-woy-project/models"
	"bayar-woy-project/responses"
	"crypto/rand"
	"fmt"
	"math/big"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

const codeExpiryMinutes = 10

func GenerateDiscordCode(c *gin.Context) {
	userID := c.GetString("userID")

	// Delete any existing verification for this user (regenerate semantics)
	config.DB.Where("user_id = ?", userID).Delete(&models.DiscordVerification{})

	code, err := generateUniqueCode()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to generate code"})
		return
	}

	expiresAt := time.Now().Add(codeExpiryMinutes * time.Minute)
	verification := models.DiscordVerification{
		UserID:    userID,
		Code:      code,
		ExpiresAt: expiresAt,
	}
	if err := config.DB.Create(&verification).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to save code"})
		return
	}

	installURL := fmt.Sprintf(
		"https://discord.com/oauth2/authorize?client_id=%s&scope=applications.commands+bot&integration_type=1",
		config.GetEnv("DISCORD_CLIENT_ID"),
	)

	c.JSON(http.StatusOK, responses.APIResponse{
		StatusCode: http.StatusOK,
		Message:    "Code generated",
		Data: responses.DiscordVerifyResponse{
			Code:        code,
			BotUsername: config.GetEnv("DISCORD_BOT_USERNAME"),
			ExpiresAt:   expiresAt,
			InstallURL:  installURL,
		},
	})
}

func generateUniqueCode() (string, error) {
	for i := 0; i < 10; i++ {
		n, err := rand.Int(rand.Reader, big.NewInt(1000000))
		if err != nil {
			return "", err
		}
		candidate := fmt.Sprintf("%06d", n.Int64())

		var existing models.DiscordVerification
		res := config.DB.Where("code = ?", candidate).First(&existing)
		if res.Error != nil {
			return candidate, nil // not found = unique
		}
	}
	return "", fmt.Errorf("failed to generate unique code after 10 attempts")
}
