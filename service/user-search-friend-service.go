package service

import (
	"bayar-woy-project/config"
	"bayar-woy-project/models"
	"bayar-woy-project/responses"

	"github.com/gin-gonic/gin"
)

func containsFriend(friendships []models.Friendship, userId string) string {
	for _, friendship := range friendships {
		if friendship.FriendID == userId {
			return "friend"
		}
	}
	return "not_friend"
}

func SearchFriend(c *gin.Context) {
	var users []models.User
	var query struct {
		Name string `json:"name" binding:"required"`
	}
	userId := c.GetString("userId")

	if err := c.ShouldBindJSON(&query); err != nil {
		c.JSON(400, gin.H{"message": "Name query parameter is required"})
		return
	}

	if err := config.DB.Where("username LIKE ?", query.Name+"%").Find(&users).Error; err != nil {
		c.JSON(500, gin.H{"message": "Failed to search users"})
		return
	}
	var friendResponses []responses.FriendResponse
	var userFriend []models.Friendship
	config.DB.Preload("Friend").Where("user_id = ?", userId).Find(&userFriend)

	for _, user := range users {
		if user.ID == userId {
			continue
		}

		friendResponses = append(friendResponses, responses.FriendResponse{
			ID:       user.ID,
			Username: user.Username,
			Status:   containsFriend(userFriend, user.ID),
		})
	}

	apiResponse := responses.APIResponse{
		StatusCode: 200,
		Message:    "Users found successfully",
		Data:       friendResponses,
	}

	c.JSON(200, apiResponse)
}

func GetAllFriends(c *gin.Context) {
	userId := c.GetString("userId")
	var friends []models.Friendship

	if err := config.DB.Where("user_id = ?", userId).Find(&friends).Error; err != nil {
		c.JSON(500, gin.H{"message": "Failed to get friends"})
		return
	}

	var friendUsers []responses.FriendResponse
	for _, friend := range friends {
		friendUsers = append(friendUsers, responses.FriendResponse{
			ID:       friend.FriendID,
			Username: friend.Friend.Username,
		})
	}

	apiResponse := responses.APIResponse{
		StatusCode: 200,
		Message:    "Friends found successfully",
		Data:       friendUsers,
	}

	c.JSON(200, apiResponse)

}