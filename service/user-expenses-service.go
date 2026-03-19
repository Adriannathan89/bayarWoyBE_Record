package service

import (
	"bayar-woy-project/config"
	"bayar-woy-project/dto"
	"bayar-woy-project/models"
	"bayar-woy-project/responses"
	"net/http"

	"github.com/gin-gonic/gin"
)

func CreateExpense(c *gin.Context) {
	userId := c.GetString("userId")
	var req dto.CreateExpenseDto

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Bad Request"})
		return
	}

	expense := models.Expense{
		Description: req.Description,
		Amount:      req.Amount,
		OwnerID:     userId,
	}

	if err := config.DB.Create(&expense).Error; err != nil {
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

func GetExpenses(c *gin.Context) {
	userId := c.GetString("userId")
	var expenses []models.Expense

	if err := config.DB.Where("owner_id = ?", userId).Find(&expenses).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to fetch expenses"})
		return
	}

	apiResponse := responses.APIResponse{
		StatusCode: http.StatusOK,
		Message:    "Expenses fetched successfully",
		Data:       expenses,
	}

	c.JSON(http.StatusOK, apiResponse)
}