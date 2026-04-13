package controller

import (
	"bayar-woy-project/guard"
	"bayar-woy-project/user_friend_service"

	"github.com/gin-gonic/gin"
)

func UserFriendController(r *gin.Engine) {
	user := r.Group("/user/friend")
	user.Use(guard.AuthMiddleware())
	{
		user.POST("/add", user_friend_service.SentFriendRequest)
		user.POST("/search", user_friend_service.SearchFriend)
		user.GET("/request", user_friend_service.GetFriendRequests)
		user.GET("", user_friend_service.GetAllFriends)
		user.PUT("/request/response", user_friend_service.FriendRequestResponse)
	}
}	