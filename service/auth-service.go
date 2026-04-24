package service

import (
	"bayar-woy-project/config"
	"bayar-woy-project/dto"
	"bayar-woy-project/guard"
	"bayar-woy-project/models"
	"bayar-woy-project/responses"

	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func Login(c *gin.Context) {
	var req dto.CreateUserDTO
	var apiResponse responses.APIResponse
	var user models.User
	var clientHost string = config.GetEnv("CLIENT_HOST")

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Bad Request"})
		return
	}

	if err := config.DB.Where("username = ?", req.Username).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid Login"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
		return
	}

	token, _ := guard.GenerateToken(user.Username, user.ID)
	refreshToken, _ := guard.GenerateRefreshToken(user.Username, user.ID)

	sesion := models.Session{
		UserID:       user.ID,
		Username:     user.Username,
		RefreshToken: refreshToken,
		IPAddress:    c.ClientIP(),
		ExpiresAt:    time.Now().Add(720 * time.Hour),
	}

	config.DB.Create(&sesion)

	c.SetSameSite(http.SameSiteNoneMode)
	c.SetCookie("token", token, 60*10, "/", clientHost, true, true)
	c.SetCookie("refresh_token", refreshToken, 3600*720, "/", clientHost, true, true)

	apiResponse = responses.APIResponse{
		StatusCode: http.StatusOK,
		Message:    "Login successful",
		Data:       user,
	}
	c.JSON(http.StatusOK, apiResponse)
}

func ValidateStillValidSession(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Session is still valid"})
}

func Logout(c *gin.Context) {
	refreshToken, _ := c.Cookie("refresh_token")
	var clientHost string = config.GetEnv("CLIENT_HOST")

	config.DB.Where("refresh_token = ?", refreshToken).Delete(&models.Session{})

	c.SetSameSite(http.SameSiteNoneMode)
	c.SetCookie("token", "", -1, "/", clientHost, true, true)
	c.SetCookie("refresh_token", "", -1, "/", clientHost, true, true)

	apiResponse := responses.APIResponse{
		StatusCode: http.StatusOK,
		Message:    "Logout successful",
		Data:       nil,
	}
	c.JSON(http.StatusOK, apiResponse)
}
