package controller

import (
	"bayar-woy-project/service"

	"github.com/gin-gonic/gin"
)

func UserController(r *gin.Engine) {
	user := r.Group("/user")
	{
		user.POST("/register", service.Register)
	}
}
