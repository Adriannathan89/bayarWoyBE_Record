package loader

import (
	"bayar-woy-project/config"
	"bayar-woy-project/models"
)

func LoadConfig() {
	config.LoadEnv()
	config.ConnectDatabase()

	// migrate all model here
	config.DB.AutoMigrate(&models.User{}, &models.Debt{}, &models.Session{})
	config.DB.AutoMigrate(&models.Friendship{}, &models.FriendRequest{})
	config.DB.AutoMigrate(&models.Record{})
}
