package controller

import (
	"bayar-woy-project/guard"
	"bayar-woy-project/service"

	"github.com/gin-gonic/gin"
)

func UserRelationController(r *gin.Engine) {
	user := r.Group("/user")
	user.Use(guard.AuthMiddleware())
	{
		user.POST("/add-friend", service.SentFriendRequest)
		user.GET("/friend-request", service.GetFriendRequests)
		user.PUT("/friend-request/response", service.FriendRequestResponse)
	}
}	