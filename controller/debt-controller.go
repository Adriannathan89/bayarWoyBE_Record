package controller

import (
	"bayar-woy-project/guard"
	"bayar-woy-project/service"

	"github.com/gin-gonic/gin"
)

func DebtController(r *gin.Engine) {
	debt := r.Group("/debt")
	debt.Use(guard.AuthMiddleware())
	{
		debt.POST("/create", service.CreateDebt)
		debt.PUT("/finish", service.FinishDebt)
	}
}