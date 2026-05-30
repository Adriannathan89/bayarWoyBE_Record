package user_profile_service_test

import (
	"bayar-woy-project/models"
	"bayar-woy-project/responses"
	"bayar-woy-project/testutil"
	ups "bayar-woy-project/user_profile_service"
	"encoding/json"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
)

func TestGenerateDiscordCodeCreatesRow(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db := testutil.SetupTestDB(t)

	os.Setenv("DISCORD_CLIENT_ID", "1234567890")
	os.Setenv("DISCORD_BOT_USERNAME", "TestBot")
	defer os.Unsetenv("DISCORD_CLIENT_ID")
	defer os.Unsetenv("DISCORD_BOT_USERNAME")

	user := models.User{Username: "alice", Password: "x"}
	db.Create(&user)

	c, rec := newJSONContext(http.MethodPost, "/user/discord/verify", nil)
	c.Set("userID", user.ID)

	ups.GenerateDiscordCode(c)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", rec.Code, rec.Body.String())
	}

	var resp struct {
		Data responses.DiscordVerifyResponse `json:"data"`
	}
	json.Unmarshal(rec.Body.Bytes(), &resp)

	if len(resp.Data.Code) != 6 {
		t.Errorf("expected 6-digit code, got %q", resp.Data.Code)
	}
	if resp.Data.BotUsername != "TestBot" {
		t.Errorf("expected botUsername 'TestBot', got %q", resp.Data.BotUsername)
	}
	if resp.Data.InstallURL == "" {
		t.Error("expected installUrl to be set")
	}
	if !resp.Data.ExpiresAt.After(time.Now().Add(8 * time.Minute)) {
		t.Errorf("expected expiresAt at least 8 minutes from now, got %v", resp.Data.ExpiresAt)
	}

	var stored models.DiscordVerification
	if err := db.Where("code = ?", resp.Data.Code).First(&stored).Error; err != nil {
		t.Errorf("expected verification row to exist: %v", err)
	}
	if stored.UserID != user.ID {
		t.Errorf("expected UserID=%s, got %s", user.ID, stored.UserID)
	}
}

func TestGenerateDiscordCodeReplacesExisting(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db := testutil.SetupTestDB(t)

	os.Setenv("DISCORD_CLIENT_ID", "1234567890")
	os.Setenv("DISCORD_BOT_USERNAME", "TestBot")
	defer os.Unsetenv("DISCORD_CLIENT_ID")
	defer os.Unsetenv("DISCORD_BOT_USERNAME")

	user := models.User{Username: "alice", Password: "x"}
	db.Create(&user)

	old := models.DiscordVerification{
		UserID:    user.ID,
		Code:      "111111",
		ExpiresAt: time.Now().Add(5 * time.Minute),
	}
	db.Create(&old)

	c, rec := newJSONContext(http.MethodPost, "/user/discord/verify", nil)
	c.Set("userID", user.ID)

	ups.GenerateDiscordCode(c)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}

	var count int64
	db.Model(&models.DiscordVerification{}).Where("user_id = ?", user.ID).Count(&count)
	if count != 1 {
		t.Errorf("expected exactly 1 verification row for user, got %d", count)
	}

	var stillOld models.DiscordVerification
	res := db.Where("code = ?", "111111").First(&stillOld)
	if res.Error == nil {
		t.Error("expected old code '111111' to be deleted")
	}
}
