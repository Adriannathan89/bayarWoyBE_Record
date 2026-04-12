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
	refreshToken, _ := guard.GenerateRefreshToken(user.ID, user.Username)

	sesion := models.Session{
		UserID:       user.ID,
		Username:     user.Username,
		RefreshToken: refreshToken,
		IPAddress:    c.ClientIP(),
		ExpiresAt:    time.Now().Add(4 * time.Hour),
	}

	config.DB.Create(&sesion)

	c.SetCookie("token", token, 60 * 15, "/", "localhost", false, true)
	c.SetCookie("refresh_token", refreshToken, 3600*4, "/", "localhost", false, true)

	apiResponse = responses.APIResponse{
		StatusCode: http.StatusOK,
		Message:    "Login successful",
		Data:       user,
	}
	c.JSON(http.StatusOK, apiResponse)
}

func ValidateStillValidSession(c *gin.Context) {
	token, err := c.Cookie("token")

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
		return
	}

	claims, err := guard.ValidateToken(token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Session is still valid", "claims": claims})
}

func GenerateNewToken(c *gin.Context) {
	refreshToken, err := c.Cookie("refresh_token")

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
		return
	}
	var sesion models.Session

	if err := config.DB.Where("refresh_token = ?", refreshToken).First(&sesion).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
		return
	}

	claims, _ := guard.GenerateToken(sesion.Username, sesion.UserID)

	c.SetCookie("token", claims, 3600, "/", "localhost", false, true)

	apiResponse := responses.APIResponse{
		StatusCode: http.StatusOK,
		Message:    "Token refreshed successfully",
		Data:       nil,
	}

	c.JSON(http.StatusOK, apiResponse)
}

func Logout(c *gin.Context) {
	refreshToken, _ := c.Cookie("refresh_token")

	config.DB.Where("refresh_token = ?", refreshToken).Delete(&models.Session{})

	c.SetCookie("token", "", -1, "/", "", false, true)
	c.SetCookie("refresh_token", "", -1, "/", "", false, true)

	apiResponse := responses.APIResponse{
		StatusCode: http.StatusOK,
		Message:    "Logout successful",
		Data:       nil,
	}
	c.JSON(http.StatusOK, apiResponse)
}
