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

	if req.FriendID == userId {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Cannot send friend request to yourself"})
		return
	}

	var existing models.Friendship
	if err := config.DB.Where("user_id = ? AND friend_id = ?", userId, req.FriendID).First(&existing).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"message": "Friend request already exists"})
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
		if err := tx.Create(&friendRequest).Error; err != nil {
			return err
		}
		if err := tx.Create(&newFriendship).Error; err != nil {
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

	if err := config.DB.Where("friend_id = ?", userId).Preload("Sender").Preload("Receiver").Find(&friendRequests).Error; err != nil {
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

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Bad Request"})
		return
	}

	if err := config.DB.Where("id = ?", req.FriendRequestID).First(&friendRequest).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Friend request not found"})
		return
	}

	if req.Action == "reject" {
		if err := config.DB.Transaction(func(tx *gorm.DB) error {
			if err := tx.Delete(&friendRequest).Error; err != nil {
				return err
			}
			if err := tx.Where("user_id = ? AND friend_id = ?", friendRequest.SenderID, friendRequest.ReceiverID).Delete(&models.Friendship{}).Error; err != nil {
				return err
			}
			return nil
		}); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to reject friend request"})
			return
		}

		c.JSON(http.StatusOK, responses.APIResponse{
			StatusCode: http.StatusOK,
			Message:    "Friend request rejected successfully",
			Data:       nil,
		})
		return
	}

	// accept
	if err := config.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Delete(&friendRequest).Error; err != nil {
			return err
		}
		if err := tx.Model(&models.Friendship{}).
			Where("user_id = ? AND friend_id = ?", friendRequest.SenderID, friendRequest.ReceiverID).
			Update("status", "accepted").Error; err != nil {
			return err
		}
		reverse := models.Friendship{
			UserID:   friendRequest.ReceiverID,
			FriendID: friendRequest.SenderID,
			Status:   "accepted",
		}
		if err := tx.Create(&reverse).Error; err != nil {
			return err
		}
		return nil
	}); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to accept friend request"})
		return
	}

	c.JSON(http.StatusOK, responses.APIResponse{
		StatusCode: http.StatusOK,
		Message:    "Friend request accepted successfully",
		Data:       nil,
	})
}
