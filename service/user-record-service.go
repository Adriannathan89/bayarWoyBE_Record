package service

import (
	"bayar-woy-project/config"
	"bayar-woy-project/dto"
	"bayar-woy-project/models"
	"bayar-woy-project/responses"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func CreateRecord(c *gin.Context) {
	userId := c.GetString("userID")
	var req dto.CreateRecordDto
	var user models.User

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Bad Request"})
		return
	}

	expense := models.Record{
		Title:       req.Title,
		Description: req.Description,
		Amount:      req.Amount,
		OwnerID:     userId,
		Type:        req.Type,
	}

	err := config.DB.Transaction(func(tx *gorm.DB) error {
		if err := config.DB.Where("id = ?", userId).First(&user).Error; err != nil {
			return err
		}

		if req.Type == "income" {
			user.Cash += req.Amount
		} else {
			user.Cash -= req.Amount
		}

		if err := config.DB.Save(&user).Error; err != nil {
			return err
		}
		if err := config.DB.Create(&expense).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to create expense"})
		return
	}

	apiResponse := responses.APIResponse{
		StatusCode: http.StatusOK,
		Message:    "Expense created successfully",
		Data:       expense,
	}

	c.JSON(http.StatusOK, apiResponse)
}