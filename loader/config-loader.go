package loader

import (
	"bayar-woy-project/bot"
	botmodel "bayar-woy-project/bot-model"
	"bayar-woy-project/config"
	"bayar-woy-project/models"
)

func LoadConfig() {
	config.LoadEnv()
	config.ConnectDatabase()
	botToken := config.GetEnv("DISBORD_BOT_TOKEN")

	if err := bot.InitDiscordBot(botToken); err != nil {
		panic("Failed to initialize Discord bot: " + err.Error())
	}

	// migrate all model here
	config.DB.AutoMigrate(&models.User{}, &models.Debt{}, &models.Session{})
	config.DB.AutoMigrate(&models.Friendship{}, &models.FriendRequest{}, &botmodel.DiscordBotOtp{})
	config.DB.AutoMigrate(&models.Record{})
}