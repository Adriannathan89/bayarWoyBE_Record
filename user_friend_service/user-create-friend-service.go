package user_friend_service

import (
	"bayar-woy-project/config"
	"bayar-woy-project/dto"
	"bayar-woy-project/models"
	"bayar-woy-project/responses"

	"net/http"

	"github.com/gin-gonic/gin"
)

func SentFriendRequest(c *gin.Context) {
	username := c.GetString("username")
	var req dto.RelationPhase1

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Bad Request"})
		return
	}

	friendRequest := models.FriendRequest{
		SenderUsername:   username,
		ReceiverUsername: req.FriendUsername,
	}

	if err := config.DB.Create(&friendRequest).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to send friend request"})
		return
	}

	apiResponse := responses.APIResponse{
		StatusCode: http.StatusOK,
		Message:    "Friend request sent successfully",
		Data:       friendRequest,
	}

	c.JSON(http.StatusOK, apiResponse)
}

func GetFriendRequests(c *gin.Context) {
	username := c.GetString("username")
	var friendRequests []models.FriendRequest
	
	if err := config.DB.Where("sender_username = ?", username).Find(&friendRequests).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to retrieve friend requests"})
		return
	}

	apiResponse := responses.APIResponse{
		StatusCode: http.StatusOK,
		Message:    "Friend requests retrieved successfully",
		Data:       friendRequests,
	}
	c.JSON(http.StatusOK, apiResponse)
}

func FriendRequestResponse(c *gin.Context) {
	var req dto.RelationPhase2
	var friendRequest models.FriendRequest

	c.ShouldBindJSON(&req)

	if req.Action == "reject" {
		config.DB.Where("id = ?", req.FriendRequestID).Delete(&friendRequest)
		
		apiResponse := responses.APIResponse{
			StatusCode: http.StatusOK,
			Message:    "Friend request rejected successfully",
			Data:       nil,
		}
		c.JSON(http.StatusOK, apiResponse)
		return
	}

	if err := config.DB.Where("id = ?", req.FriendRequestID).First(&friendRequest).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Friend request not found"})
		return
	}

	config.DB.Where("id = ?", req.FriendRequestID).Delete(&friendRequest)
	

	newFriendship := models.Friendship{
		UserID:   friendRequest.Sender.ID,
		FriendID: friendRequest.Receiver.ID,
	}

	config.DB.Create(&newFriendship)

	apiResponse := responses.APIResponse{
		StatusCode: http.StatusOK,
		Message:    "Friend request accepted successfully",
		Data:       newFriendship,
	}

	c.JSON(http.StatusOK, apiResponse)
}