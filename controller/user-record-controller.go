package controller

import (
	"bayar-woy-project/guard"
	"bayar-woy-project/service"

	"github.com/gin-gonic/gin"
)

func UserRecordController(r *gin.Engine) {
	user := r.Group("/user")
	user.Use(guard.AuthMiddleware())
	{
		user.POST("/record", service.CreateRecord)
	}
}