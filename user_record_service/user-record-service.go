package user_record_service

import (
	"bayar-woy-project/config"
	"bayar-woy-project/dto"
	"bayar-woy-project/models"
	"bayar-woy-project/responses"
	"bayar-woy-project/slm"
	"net/http"
	"time"

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

	parsedTime, err := time.Parse("2006-01-02", req.Date)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid date format"})
		return
	}

	slmResult := slm.Classify(req.Title)
	recordType := "expense"
	if slmResult.TransactionType == "pemasukan" {
		recordType = "income"
	}

	expense := models.Record{
		Title:       req.Title,
		Description: req.Description,
		Amount:      req.Amount,
		Category:    slmResult.Category,
		OwnerID:     userId,
		Type:        recordType,
		CreatedAt:   parsedTime,
	}

	err2 := config.DB.Transaction(func(tx *gorm.DB) error {
		if err := config.DB.Where("id = ?", userId).First(&user).Error; err != nil {
			return err
		}

		if recordType == "income" {
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

	if err2 != nil {
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
