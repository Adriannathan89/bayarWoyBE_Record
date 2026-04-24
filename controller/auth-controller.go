package controller

import (
	"bayar-woy-project/guard"
	"bayar-woy-project/service"

	"github.com/gin-gonic/gin"
)

func AuthController(r *gin.Engine) {
	auth := r.Group("/auth")
	{
		auth.POST("/login", service.Login)
		auth.POST("/logout", service.Logout)
	}
	auth.Use(guard.AuthMiddleware())
	{
		auth.GET("/validate-session", service.ValidateStillValidSession)
	}
}
