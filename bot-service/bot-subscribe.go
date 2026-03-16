package botservice

import (
	botmodel "bayar-woy-project/bot-model"
	"bayar-woy-project/config"
	"bayar-woy-project/models"
	"math/rand"
	"time"

	"fmt"
)

func RegisterUserToBot(username string, discordID string) string {
	var user models.User
	if err := config.DB.Where("username = ?", username).First(&user).Error; err != nil {
		return "User not found"
	}
	user.DiscordID = &discordID
	if err := config.DB.Save(&user).Error; err != nil {
		return "Failed to update user with Discord ID"
	}

	code := fmt.Sprintf("%06d", rand.Intn(1000000))
	otp := botmodel.DiscordBotOtp{
		UserID:    user.ID,
		OTP:       code,
		ExpiredAt: time.Now().Add(5 * time.Minute),
	}
	
	if err := config.DB.Create(&otp).Error; err != nil {
		return "Failed to create OTP"
	}

	res := "otp code: " + code + "\nThis code will expire in 5 minutes."

	return res
}