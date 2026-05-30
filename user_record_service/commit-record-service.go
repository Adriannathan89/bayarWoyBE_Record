package user_record_service

import (
	"bayar-woy-project/config"
	"bayar-woy-project/discord"
	"bayar-woy-project/dto"
	"bayar-woy-project/models"
	"bayar-woy-project/responses"
	"bayar-woy-project/slm"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

var incomeCategories = map[string]bool{"gaji": true, "hadiah": true}

func CommitRecord(c *gin.Context) {
	userID := c.GetString("userID")
	var req dto.CommitRecordDto

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Bad Request"})
		return
	}

	var record models.Record
	if err := config.DB.Preload("Categories").
		Where("id = ? AND owner_id = ?", req.RecordID, userID).
		First(&record).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Record not found"})
		return
	}

	if record.IsCommitted {
		c.JSON(http.StatusConflict, gin.H{"message": "Record already committed"})
		return
	}

	// Resolve stored primary and secondary categories
	storedPrimary := ""
	storedSecondary := ""
	for _, cat := range record.Categories {
		if cat.Type == "primary" {
			storedPrimary = cat.Name
		} else if cat.Type == "secondary" {
			storedSecondary = cat.Name
		}
	}

	// Final category: from request if provided, otherwise stored
	finalCategory := req.Category
	if finalCategory == "" {
		finalCategory = storedPrimary
	}

	// Final secondary: from request if provided, otherwise same-name as new primary if primary changed, otherwise stored
	finalSecondary := storedSecondary
	if req.SecondaryCategory != "" {
		finalSecondary = req.SecondaryCategory
	} else if req.Category != "" && req.Category != storedPrimary {
		finalSecondary = finalCategory
	}

	// Derive type from final category
	finalType := "expense"
	if incomeCategories[finalCategory] {
		finalType = "income"
	}

	err := config.DB.Transaction(func(tx *gorm.DB) error {
		// Update categories and type if primary changed
		if req.Category != "" && req.Category != storedPrimary {
			var newPrimary, newSecondary models.Category
			tx.Where("name = ? AND type = 'primary'", finalCategory).First(&newPrimary)
			tx.Where("name = ? AND type = 'secondary'", finalSecondary).First(&newSecondary)

			var newCats []models.Category
			if newPrimary.ID != "" {
				newCats = append(newCats, newPrimary)
			}
			if newSecondary.ID != "" {
				newCats = append(newCats, newSecondary)
			}
			if err := tx.Model(&record).Association("Categories").Replace(newCats); err != nil {
				return err
			}
			record.Type = finalType
			if err := tx.Save(&record).Error; err != nil {
				return err
			}
		}

		// Update user cash
		var user models.User
		if err := tx.Where("id = ?", userID).First(&user).Error; err != nil {
			return err
		}
		if finalType == "income" {
			user.Cash += record.Amount
		} else {
			user.Cash -= record.Amount
		}
		if err := tx.Save(&user).Error; err != nil {
			return err
		}

		// Mark committed
		if err := tx.Model(&record).Update("is_committed", true).Error; err != nil {
			return err
		}
		record.IsCommitted = true
		return nil
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to commit record"})
		return
	}

	// Send SLM feedback (best-effort, outside transaction)
	slm.Feedback(record.Title, finalCategory, finalSecondary)

	// Send Discord notification (best-effort, fire-and-forget)
	var fullUser models.User
	if err := config.DB.First(&fullUser, "id = ?", userID).Error; err == nil {
		// Re-load record with categories for notification
		var notifRecord models.Record
		if err := config.DB.Preload("Categories").First(&notifRecord, "id = ?", record.ID).Error; err == nil {
			go discord.NotifyCommit(fullUser, notifRecord)
		}
	}

	c.JSON(http.StatusOK, responses.APIResponse{
		StatusCode: http.StatusOK,
		Message:    "Record committed",
		Data:       record,
	})
}
