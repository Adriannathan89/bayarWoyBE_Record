package main

import (
	"bayar-woy-project/config"
	"bayar-woy-project/controller"
	"bayar-woy-project/loader"
	"strings"
	"time"

	"github.com/gin-contrib/cors"

	"github.com/gin-gonic/gin"
)

func parseAllowedOrigins(raw string) []string {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return []string{}
	}

	raw = strings.TrimPrefix(raw, "[")
	raw = strings.TrimSuffix(raw, "]")

	parts := strings.Split(raw, ",")
	origins := make([]string, 0, len(parts))
	for _, part := range parts {
		origin := strings.TrimSpace(part)
		origin = strings.Trim(origin, "\"")
		if origin != "" {
			origins = append(origins, origin)
		}
	}

	return origins
}

func main() {
	loader.LoadConfig()

	r := gin.Default()
	allowedOrigins := parseAllowedOrigins(config.GetEnv("ALLOWED_ORIGINS"))

	r.Use(cors.New(cors.Config{
		AllowOrigins:     allowedOrigins,
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	controller.AuthController(r)
	controller.DebtController(r)
	controller.UserController(r)
	controller.UserFriendController(r)
	controller.UserRecordController(r)

	r.Run(":" + config.GetEnv("PORT"))
}
