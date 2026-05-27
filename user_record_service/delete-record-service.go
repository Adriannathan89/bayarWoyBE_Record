package user_record_service

import (
	"bayar-woy-project/config"
	"bayar-woy-project/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func DeleteRecord(c *gin.Context) {
	userID := c.GetString("userID")
	recordID := c.Param("id")

	var record models.Record
	if err := config.DB.Where("id = ? AND owner_id = ?", recordID, userID).First(&record).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Record not found"})
		return
	}

	if record.IsCommitted {
		c.JSON(http.StatusConflict, gin.H{"message": "Cannot delete a committed record"})
		return
	}

	if err := config.DB.Model(&record).Association("Categories").Clear(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to clear categories"})
		return
	}

	if err := config.DB.Delete(&record).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to delete record"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Record deleted"})
}
