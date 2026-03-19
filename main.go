package main

import (
	"bayar-woy-project/config"
	"bayar-woy-project/controller"
	"bayar-woy-project/loader"
	"time"

	"github.com/gin-contrib/cors"

	"github.com/gin-gonic/gin"
)

func main() {
	loader.LoadConfig()
	r := gin.Default()
	allowedOrigins := config.GetEnv("ALLOWED_ORIGINS")

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{string(allowedOrigins)},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: true,
		MaxAge: 12 * time.Hour,
	}))

	controller.AuthController(r)
	controller.UserController(r)
	controller.UserRelationController(r)
	controller.TransactionController(r)
	controller.UserExpenseController(r)

	r.Run(":" + config.GetEnv("PORT"))
}