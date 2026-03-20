package controller

import (
	"bayar-woy-project/guard"
	"bayar-woy-project/service"

	"github.com/gin-gonic/gin"
)

func UserFriendController(r *gin.Engine) {
	user := r.Group("/user")
	user.Use(guard.AuthMiddleware())
	{
		user.POST("/add-friend", service.SentFriendRequest)
		user.POST("/friend/search", service.SearchFriend)
		user.GET("/friend-request", service.GetFriendRequests)
		user.GET("/friend", service.GetAllFriends)
		user.PUT("/friend-request/response", service.FriendRequestResponse)
	}
}	