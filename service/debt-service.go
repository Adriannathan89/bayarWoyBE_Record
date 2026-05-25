package service

import (
	"bayar-woy-project/config"
	"bayar-woy-project/dto"
	"bayar-woy-project/models"
	"bayar-woy-project/responses"
	"bayar-woy-project/slm"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

var (
	errForbidden        = fmt.Errorf("forbidden")
	errAlreadyCompleted = fmt.Errorf("already_completed")
)

func CreateDebt(c *gin.Context) {
	var req dto.DebtDTO
	var owner models.User
	var debtor models.User
	userId := c.GetString("userID")

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Bad Request"})
		return
	}

	category := slm.ClassifyTitle(req.Description)

	transactionModel := models.Debt{
		Amount:      req.Amount,
		Description: req.Description,
		Category:    category,
		DebtorID:    req.DebtorID,
		OwnerID:     userId,
		Status:      "pending",
	}

	err := config.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("id = ?", userId).First(&owner).Error; err != nil {
			return err
		}
		if err := tx.Where("id = ?", req.DebtorID).First(&debtor).Error; err != nil {
			return err
		}

		owner.Receivable += req.Amount
		owner.Cash -= req.Amount
		debtor.Debt += req.Amount

		if err := tx.Save(&owner).Error; err != nil {
			return err
		}
		if err := tx.Save(&debtor).Error; err != nil {
			return err
		}
		if err := tx.Create(&transactionModel).Error; err != nil {
			return err
		}

		ownerRecord := models.Record{
			Title:       "Hutang ke " + debtor.Username,
			Description: req.Description,
			Category:    category,
			Amount:      req.Amount,
			OwnerID:     userId,
			Type:        "expense",
			CreatedAt:   time.Now(),
		}
		if err := tx.Create(&ownerRecord).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to create debt"})
		return
	}

	c.JSON(http.StatusOK, responses.APIResponse{
		StatusCode: http.StatusOK,
		Message:    "Debt created successfully",
		Data:       transactionModel,
	})
}

func FinishDebt(c *gin.Context) {
	var req dto.FinishDebtDTO
	userId := c.GetString("userID")

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Bad Request"})
		return
	}

	var transaction models.Debt

	err := config.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Preload("Owner").Preload("Debtor").
			Where("id = ?", req.DebtID).First(&transaction).Error; err != nil {
			return err
		}

		if transaction.DebtorID != userId {
			return errForbidden
		}

		if transaction.Status != "pending" {
			return errAlreadyCompleted
		}

		owner := transaction.Owner
		debtor := transaction.Debtor

		owner.Receivable -= transaction.Amount
		owner.Cash += transaction.Amount
		debtor.Debt -= transaction.Amount
		debtor.Cash -= transaction.Amount
		transaction.Status = "completed"

		if err := tx.Save(&owner).Error; err != nil {
			return err
		}
		if err := tx.Save(&debtor).Error; err != nil {
			return err
		}
		if err := tx.Save(&transaction).Error; err != nil {
			return err
		}

		now := time.Now()
		ownerRecord := models.Record{
			Title:     "Piutang lunas dari " + debtor.Username,
			Amount:    transaction.Amount,
			OwnerID:   owner.ID,
			Type:      "income",
			CreatedAt: now,
		}
		if err := tx.Create(&ownerRecord).Error; err != nil {
			return err
		}

		debtorRecord := models.Record{
			Title:     "Bayar hutang ke " + owner.Username,
			Amount:    transaction.Amount,
			OwnerID:   debtor.ID,
			Type:      "expense",
			CreatedAt: now,
		}
		if err := tx.Create(&debtorRecord).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		if errors.Is(err, errForbidden) {
			c.JSON(http.StatusForbidden, gin.H{"message": "Only the debtor can finish a debt"})
			return
		}
		if errors.Is(err, errAlreadyCompleted) {
			c.JSON(http.StatusConflict, gin.H{"message": "Debt is already completed"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to finish debt"})
		return
	}

	c.JSON(http.StatusOK, responses.APIResponse{
		StatusCode: http.StatusOK,
		Message:    "Debt finished successfully",
		Data:       transaction,
	})
}

func LoadAllDebt(c *gin.Context) {
	userId := c.GetString("userID")
	var user models.User

	if err := config.DB.Preload("Debts.Owner").Preload("Debts.Debtor").Where("id = ?", userId).First(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to load debts"})
		return
	}

	apiResponse := responses.APIResponse{
		StatusCode: http.StatusOK,
		Message:    "Debts loaded successfully",
		Data: gin.H{
			"debts": user.Debts,
		},
	}
	c.JSON(http.StatusOK, apiResponse)
}

func LoadOwedDebt(c *gin.Context) {
	userId := c.GetString("userID")
	var debts []models.Debt

	if err := config.DB.Preload("Owner").Preload("Debtor").
		Where("debtor_id = ? AND status = ?", userId, "pending").
		Find(&debts).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to load owed debts"})
		return
	}

	c.JSON(http.StatusOK, responses.APIResponse{
		StatusCode: http.StatusOK,
		Message:    "Owed debts loaded successfully",
		Data:       gin.H{"debts": debts},
	})
}