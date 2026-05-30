package controller

import (
	"bayar-woy-project/guard"
	"bayar-woy-project/user_profile_service"

	"github.com/gin-gonic/gin"
)

func UserProfileController(r *gin.Engine) {
	user := r.Group("/user")
	user.Use(guard.AuthMiddleware())
	{
		user.GET("/profile", user_profile_service.GetProfile)
		user.PUT("/profile", user_profile_service.UpdateProfile)
		user.POST("/discord/verify", user_profile_service.GenerateDiscordCode)
		user.GET("/discord/status", user_profile_service.GetDiscordStatus)
		user.DELETE("/discord", user_profile_service.DisconnectDiscord)
		user.PUT("/discord/notification", user_profile_service.UpdateNotifSettings)
	}
}
