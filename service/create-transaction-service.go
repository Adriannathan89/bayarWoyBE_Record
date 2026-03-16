package service

import (
	"bayar-woy-project/config"
	"bayar-woy-project/dto"
	"bayar-woy-project/models"
	"bayar-woy-project/responses"

	"github.com/gin-gonic/gin"

	"net/http"
)

func CreateTransaction(c *gin.Context) {
	var req dto.TransactionDTO
	var owner models.User
	var debtor models.User

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Bad Request"})
		return
	}

	transactionModel := models.Transaction{
		Amount:      req.Amount,
		Description: req.Description,
		DebtorID:    req.DebtorID,
		OwnerID:     req.OwnerID,
		Status:      "pending",
	}

	if config.DB.Create(&transactionModel).Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Failed to create transaction"})
		return
	}

	config.DB.Where("id = ?", req.OwnerID).First(&owner)
	config.DB.Where("id = ?", req.DebtorID).First(&debtor)

	owner.Receivable += req.Amount
	debtor.Debt += req.Amount

	config.DB.Save(&owner)
	config.DB.Save(&debtor)

	apiResponse := responses.APIResponse{
		StatusCode: http.StatusOK,
		Message:    "Transaction created successfully",
		Data:       transactionModel,
	}

	c.JSON(http.StatusOK, apiResponse)
}

func GetTransactions(c *gin.Context) {
	username := c.GetString("username")
	var user models.User
	if err := config.DB.Where("username = ?", username).First(&user).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "User not found"})
		return
	}
	var transactions []models.Transaction
	if err := config.DB.Where("owner_id = ?", user.ID).Preload("Debtor").Find(&transactions).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Failed to retrieve transactions"})
		return
	}

	apiResponse := responses.APIResponse{
		StatusCode: http.StatusOK,
		Message:    "Transactions retrieved successfully",
		Data:       transactions,
	}

	c.JSON(http.StatusOK, apiResponse)
}


func FinishTransaction(c *gin.Context) {
	var req dto.FinishTransactionDTO

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Bad Request"})
		return
	}

	var transaction models.Transaction
	if err := config.DB.Preload("Owner").Where("id = ?", req.TransactionID).First(&transaction).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Transaction not found"})
		return
	}
	owner := transaction.Owner

	owner.Receivable -= transaction.Amount
	owner.Cash += transaction.Amount
	transaction.Status = "completed"

	config.DB.Save(&owner)
	config.DB.Save(&transaction)
}