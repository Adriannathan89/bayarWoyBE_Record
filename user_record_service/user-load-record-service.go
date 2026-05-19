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
	var expenses []responses.RecordResponse
	var incomes []responses.RecordResponse
	var debts []responses.RecordResponse

	if err := config.DB.Preload("Records").Preload("Debts").Where("id = ?", userId).First(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to load transactions"})
		return
	}

	for _, record := range user.Records {
		if record.Type == "expense" {
			expenses = append(expenses, responses.RecordResponse{
				ID:          record.ID,
				Title:       record.Title,
				Description: record.Description,
				Amount:      record.Amount,
				Type:        record.Type,
				CreatedAt:   record.CreatedAt.Format("2006-01-02 15:04:05"),
			})
		} else if record.Type == "income" {
			incomes = append(incomes, responses.RecordResponse{
				ID:          record.ID,
				Title:       record.Title,
				Description: record.Description,
				Amount:      record.Amount,
				Type:        record.Type,
				CreatedAt:   record.CreatedAt.Format("2006-01-02 15:04:05"),
			})
		}
	}

	for _, debt := range user.Debts {
		debts = append(debts, responses.RecordResponse{
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
			"expenses":   expenses,
			"incomes":    incomes,
			"debts":      debts,
			"cash":       user.Cash,
			"debt":       user.Debt,
			"receivable": user.Receivable,
			"balance":    user.Cash + user.Receivable - user.Debt,
		},
	}
	c.JSON(http.StatusOK, apiResponse)
}
