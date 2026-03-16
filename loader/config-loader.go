package loader

import (
	"bayar-woy-project/config"
	"bayar-woy-project/models"
)

func LoadConfig() {
	config.LoadEnv()
	config.ConnectDatabase()


	// migrate all model here
	config.DB.AutoMigrate(&models.User{}, &models.Transaction{}, &models.Sesion{})
	config.DB.AutoMigrate(&models.Friendship{}, &models.FriendRequest{})
}