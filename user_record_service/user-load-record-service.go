package user_record_service

import (
	"bayar-woy-project/config"
	"bayar-woy-project/models"
	"bayar-woy-project/responses"
	"net/http"

	"github.com/gin-gonic/gin"
)

func LoadAllRecords(c *gin.Context) {
	userId := c.GetString("userID")
	var user models.User
	var records []responses.RecordResponse
	

	if err := config.DB.Preload("Records", "Debts").Where("id = ?", userId).First(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to load transactions"})
		return
	}

	for _, record := range user.Records {
		records = append(records, responses.RecordResponse{
			ID:          record.ID,
			Title:       record.Title,
			Description: record.Description,
			Amount:      record.Amount,
			Type:        record.Type,
		})
	}

	for _, debt := range user.Debts {
		records = append(records, responses.RecordResponse{
			ID:          debt.ID,
			Title:       "Debt with " + debt.Debtor.Username,
			Description: debt.Description,
			Amount:      debt.Amount,
			Type:        "debt",
		})
	}

	apiResponse := responses.APIResponse{
		StatusCode: http.StatusOK,
		Message:    "Transactions loaded successfully",
		Data: gin.H{
			"records": records,
			"cash":    user.Cash,
			"debt":    user.Debt,
			"receivable": user.Receivable,
		},
	}
	c.JSON(http.StatusOK, apiResponse)
}