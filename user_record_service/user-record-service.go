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
)

func CreateRecord(c *gin.Context) {
	userId := c.GetString("userID")
	var req dto.CreateRecordDto

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

	var categories []models.Category
	if slmResult.Category != "" {
		var primary models.Category
		config.DB.Where("name = ? AND type = 'primary'", slmResult.Category).First(&primary)
		if primary.ID != "" {
			categories = append(categories, primary)
		}
	}
	if slmResult.SecondaryCategory != "" {
		var secondary models.Category
		config.DB.Where("name = ? AND type = 'secondary'", slmResult.SecondaryCategory).First(&secondary)
		if secondary.ID != "" {
			categories = append(categories, secondary)
		}
	}

	record := models.Record{
		Title:       req.Title,
		Description: req.Description,
		Amount:      req.Amount,
		Categories:  categories,
		OwnerID:     userId,
		Type:        recordType,
		CreatedAt:   parsedTime,
		IsCommitted: false,
	}

	if err := config.DB.Create(&record).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to create record"})
		return
	}

	apiResponse := responses.APIResponse{
		StatusCode: http.StatusOK,
		Message:    "Record created successfully",
		Data:       record,
	}
	c.JSON(http.StatusOK, apiResponse)
}
