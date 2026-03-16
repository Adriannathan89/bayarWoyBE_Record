package main

import (
	"bayar-woy-project/config"
	"bayar-woy-project/controller"
	"bayar-woy-project/loader"

	"github.com/gin-gonic/gin"
)

func main() {
	loader.LoadConfig()
	r := gin.Default()

	controller.AuthController(r)
	controller.UserController(r)
	controller.UserRelationController(r)
	controller.TransactionController(r)

	r.Run(":" + config.GetEnv("PORT"))
}