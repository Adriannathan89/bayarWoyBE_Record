package controller

import (
	"bayar-woy-project/guard"
	"bayar-woy-project/service"

	"github.com/gin-gonic/gin"
)

func TransactionController(r *gin.Engine) {
	transaction := r.Group("/transaction")
	transaction.Use(guard.AuthMiddleware())
	{
		transaction.POST("/create", service.CreateTransaction)
		transaction.GET("/list", service.GetTransactions)
		transaction.PUT("/finish", service.FinishTransaction)
	}
}