package user_friend_service

import (
	"bayar-woy-project/config"
	"bayar-woy-project/dto"
	"bayar-woy-project/models"
	"bayar-woy-project/responses"
	"net/http"

	"gorm.io/gorm"

	"github.com/gin-gonic/gin"
)

func SentFriendRequest(c *gin.Context) {
	userId := c.GetString("userID")
	var req dto.RelationPhase1

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Bad Request"})
		return
	}

	friendRequest := models.FriendRequest{
		SenderID:   userId,
		ReceiverID: req.FriendID,
	}
	newFriendship := models.Friendship{
		UserID:   userId,
		FriendID: req.FriendID,
		Status:   "pending",
	}

	if err := config.DB.Transaction(func(tx *gorm.DB) error {
		if err := config.DB.Create(&friendRequest).Error; err != nil {
			return err
		}
		if err := config.DB.Create(&newFriendship).Error; err != nil {
			return err
		}
		return nil
	}); err != nil {
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
	userId := c.GetString("userID")
	var friendRequests []models.FriendRequest
	
	if err := config.DB.Where("sender_id = ?", userId).Find(&friendRequests).Error; err != nil {
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
	var friendship models.Friendship

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
	config.DB.Model(&friendship).Where("user_id = ? AND friend_id = ?", friendRequest.SenderID, friendRequest.ReceiverID).Update("status", "accepted")
	

	apiResponse := responses.APIResponse{
		StatusCode: http.StatusOK,
		Message:    "Friend request accepted successfully",
		Data:       friendship,
	}

	c.JSON(http.StatusOK, apiResponse)
}