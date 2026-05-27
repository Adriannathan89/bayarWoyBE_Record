package controller

import (
	"bayar-woy-project/guard"
	"bayar-woy-project/user_record_service"

	"github.com/gin-gonic/gin"
)

func UserRecordController(r *gin.Engine) {
	user := r.Group("/user")
	user.Use(guard.AuthMiddleware())
	{
		user.POST("/record", user_record_service.CreateRecord)
		user.GET("/records", user_record_service.LoadAllRecords)
		user.PUT("/record/commit", user_record_service.CommitRecord)
		user.DELETE("/record/:id", user_record_service.DeleteRecord)
	}
}