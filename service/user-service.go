package service

import (
	"bayar-woy-project/config"
	"bayar-woy-project/dto"
	"bayar-woy-project/models"
	"bayar-woy-project/responses"

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
