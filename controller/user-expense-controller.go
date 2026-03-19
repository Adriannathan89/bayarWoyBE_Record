package controller

import (
	"bayar-woy-project/guard"
	"bayar-woy-project/service"

	"github.com/gin-gonic/gin"
)

func UserExpenseController(r *gin.Engine) {
	user := r.Group("/user")
	user.Use(guard.AuthMiddleware())
	{
		user.POST("/expenses", service.CreateExpense)
		user.GET("/expenses", service.GetExpenses)
	}
}