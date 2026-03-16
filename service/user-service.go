package service

import (
	botmodel "bayar-woy-project/bot-model"
	"bayar-woy-project/config"
	"bayar-woy-project/dto"
	"bayar-woy-project/models"
	"bayar-woy-project/responses"
	"time"

	"github.com/gin-gonic/gin"

	"net/http"

	"golang.org/x/crypto/bcrypt"
)

func Register(c *gin.Context) {
	var req dto.CreateUserDTO
	var apiResponse responses.APIResponse

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Bad Request"})
		return
	}

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)

	user := models.User{
		Username:   req.Username,
		Password:   string(hashedPassword),
		Debt:       0,
		Cash:       0,
		Receivable: 0,
	}

	result := config.DB.Create(&user)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to create user"})
		return
	}

	apiResponse = responses.APIResponse{
		StatusCode: http.StatusCreated,
		Message:    "User created successfully",
		Data:       user,
	}

	c.JSON(http.StatusCreated, apiResponse)
}

func ValidateOtp(c *gin.Context) {
	var req dto.ValidateOtpDto
	var otp botmodel.DiscordBotOtp

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Bad Request"})
		return
	}

	if err := config.DB.Where("user_id = ?", req.UserID).First(&otp).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "OTP not found"})
		return
	}

	if otp.OTP != req.OTP || !otp.ExpiredAt.After(time.Now()) {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid OTP"})
		return
	}

	apiResponse := responses.APIResponse{
		StatusCode: http.StatusOK,
		Message:    "OTP validated successfully",
		Data:       nil,
	}
	c.JSON(http.StatusOK, apiResponse)
}
