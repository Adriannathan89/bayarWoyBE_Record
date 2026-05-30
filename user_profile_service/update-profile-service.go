package user_profile_service

import (
	"bayar-woy-project/config"
	"bayar-woy-project/dto"
	"bayar-woy-project/models"
	"bayar-woy-project/responses"
	"net/http"

	"github.com/gin-gonic/gin"
)

func UpdateProfile(c *gin.Context) {
	userID := c.GetString("userID")

	var req dto.UpdateProfileDto
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Bad Request"})
		return
	}

	var existing models.User
	res := config.DB.Where("username = ? AND id != ?", req.Username, userID).First(&existing)
	if res.Error == nil {
		c.JSON(http.StatusConflict, gin.H{"message": "Username already taken"})
		return
	}

	if err := config.DB.Model(&models.User{}).Where("id = ?", userID).
		Update("username", req.Username).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to update profile"})
		return
	}

	c.JSON(http.StatusOK, responses.APIResponse{
		StatusCode: http.StatusOK,
		Message:    "Profile updated",
		Data:       gin.H{"success": true},
	})
}
