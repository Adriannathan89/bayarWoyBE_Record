package testutil

import (
	botmodel "bayar-woy-project/bot-model"
	"bayar-woy-project/config"
	"bayar-woy-project/models"
	"fmt"
	"testing"
	"time"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

func SetupTestDB(t *testing.T) *gorm.DB {
	t.Helper()

	dsn := fmt.Sprintf("file:testdb_%d?mode=memory&cache=shared", time.Now().UnixNano())
	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open sqlite test db: %v", err)
	}

	err = db.AutoMigrate(
		&models.User{},
		&models.Session{},
		&models.Record{},
		&models.Debt{},
		&models.FriendRequest{},
		&models.Friendship{},
		&botmodel.DiscordBotOtp{},
	)
	if err != nil {
		t.Fatalf("failed to migrate sqlite test db: %v", err)
	}

	config.DB = db
	return db
}
