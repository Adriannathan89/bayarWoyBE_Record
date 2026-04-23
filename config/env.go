package config

import (
	"log"
	"os"
	"sync"

	"github.com/joho/godotenv"
)

var loadEnvOnce sync.Once

func LoadEnv() {
	loadEnvOnce.Do(func() {
		err := godotenv.Load()
		if err != nil && !os.IsNotExist(err) {
			log.Printf("Error loading .env file: %v", err)
		}
	})
}

func GetEnv(key string) string {
	LoadEnv()
	return os.Getenv(key)
}
