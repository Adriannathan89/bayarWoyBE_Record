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

func CreateDebt(c *gin.Context) {
	var req dto.DebtDTO
	var owner models.User
	var debtor models.User
	userId := c.GetString("userId")

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Bad Request"})
		return
	}

	transactionModel := models.Debt{
		Amount:      req.Amount,
		Description: req.Description,
		DebtorID:    req.DebtorID,
		OwnerID:     userId,
		Status:      "pending",
	}

	err := config.DB.Transaction(func(tx *gorm.DB) error {
		config.DB.Where("id = ?", userId).First(&owner)
		config.DB.Where("id = ?", req.DebtorID).First(&debtor)

		owner.Receivable += req.Amount
		debtor.Debt += req.Amount

		if err := config.DB.Save(&owner).Error; err != nil {
			return err
		}
		if err := config.DB.Save(&debtor).Error; err != nil {
			return err
		}
		if err := config.DB.Create(&transactionModel).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to create debt"})
		return
	}

	apiResponse := responses.APIResponse{
		StatusCode: http.StatusOK,
		Message:    "Debt created successfully",
		Data:       transactionModel,
	}

	c.JSON(http.StatusOK, apiResponse)
}

func FinishDebt(c *gin.Context) {
	var req dto.FinishDebtDTO

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Bad Request"})
		return
	}

	var transaction models.Debt

	err := config.DB.Transaction(func(tx *gorm.DB) error {
		if err := config.DB.Preload("Owner").Where("id = ?", req.DebtID).First(&transaction).Error; err != nil {
			return err
		}
		owner := transaction.Owner

		owner.Receivable -= transaction.Amount
		owner.Cash += transaction.Amount
		transaction.Status = "completed"

		config.DB.Save(&owner)
		config.DB.Save(&transaction)
		return nil
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to finish debt"})
		return
	}

	apiResponse := responses.APIResponse{
		StatusCode: http.StatusOK,
		Message:    "Debt finished successfully",
		Data:       transaction,
	}

	c.JSON(http.StatusOK, apiResponse)
}
