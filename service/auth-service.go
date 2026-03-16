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

	token, _ := guard.GenerateToken(user.ID, user.Username)
	refreshToken, _ := guard.GenerateRefreshToken(user.ID, user.Username)

	sesion := models.Sesion{
		UserID: user.ID,
		Username: user.Username,
		RefreshToken: refreshToken,
		IPAddress: c.ClientIP(),
		ExpiresAt: time.Now().Add(24 * time.Hour),
	}

	config.DB.Create(&sesion)

	c.SetCookie("token", token, 3600, "/", "localhost", false, true)
	c.SetCookie("refresh_token", refreshToken, 3600*24, "/", "localhost", false, true)


	apiResponse = responses.APIResponse{
		StatusCode: http.StatusOK,
		Message:    "Login successful",
		Data:       user,
	}
	c.JSON(http.StatusOK, apiResponse)
}

func GenerateNewToken(c *gin.Context) {
	refreshToken, _ := c.Cookie("refresh_token")
	var sesion models.Sesion

	if err := config.DB.Where("refresh_token = ?", refreshToken).First(&sesion).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
		return
	}

	claims, _ := guard.GenerateToken(sesion.UserID, sesion.Username)

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

	config.DB.Where("refresh_token = ?", refreshToken).Delete(&models.Sesion{})

	c.SetCookie("token", "", -1, "/", "", false, true)
	c.SetCookie("refresh_token", "", -1, "/", "", false, true)

	apiResponse := responses.APIResponse{
		StatusCode: http.StatusOK,
		Message:    "Logout successful",
		Data:       nil,
	}
	c.JSON(http.StatusOK, apiResponse)
}