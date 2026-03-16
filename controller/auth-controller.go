package controller

import (
	"bayar-woy-project/service"

	"github.com/gin-gonic/gin"
)

func AuthController(r *gin.Engine) {
	auth := r.Group("/auth")
	{
		auth.POST("/login", service.Login)
		auth.POST("/refresh", service.GenerateNewToken)
		auth.POST("/logout", service.Logout)
	}
}