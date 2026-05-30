package user_profile_service_test

import (
	"bayar-woy-project/models"
	"bayar-woy-project/responses"
	"bayar-woy-project/testutil"
	ups "bayar-woy-project/user_profile_service"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestGetProfileReturnsNotConnectedState(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db := testutil.SetupTestDB(t)

	user := models.User{Username: "alice", Password: "x"}
	db.Create(&user)

	c, rec := newJSONContext(http.MethodGet, "/user/profile", nil)
	c.Set("userID", user.ID)

	ups.GetProfile(c)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", rec.Code, rec.Body.String())
	}

	var resp struct {
		Data responses.ProfileResponse `json:"data"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to parse: %v", err)
	}

	if resp.Data.Username != "alice" {
		t.Errorf("expected username 'alice', got %q", resp.Data.Username)
	}
	if resp.Data.Discord.Connected {
		t.Error("expected discord.connected=false")
	}
}

func TestGetProfileReturnsConnectedState(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db := testutil.SetupTestDB(t)

	discordID := "999000111"
	discordUsername := "alicediscord"
	user := models.User{
		Username:                  "alice",
		Password:                  "x",
		DiscordID:                 &discordID,
		DiscordUsername:           &discordUsername,
		IsValidated:               true,
		DiscordCommitNotifEnabled: true,
		DiscordWeeklyNotifEnabled: false,
	}
	db.Create(&user)
	// GORM treats Go zero-value bool (false) as "not set" when column has default:true,
	// so explicitly force the WeeklyNotifEnabled field to false via Update.
	db.Model(&models.User{}).Where("id = ?", user.ID).
		Update("discord_weekly_notif_enabled", false)

	c, rec := newJSONContext(http.MethodGet, "/user/profile", nil)
	c.Set("userID", user.ID)

	ups.GetProfile(c)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}

	var resp struct {
		Data responses.ProfileResponse `json:"data"`
	}
	json.Unmarshal(rec.Body.Bytes(), &resp)

	if !resp.Data.Discord.Connected {
		t.Error("expected discord.connected=true")
	}
	if resp.Data.Discord.Username != "alicediscord" {
		t.Errorf("expected discord.username='alicediscord', got %q", resp.Data.Discord.Username)
	}
	if !resp.Data.Discord.CommitNotifEnabled {
		t.Error("expected commitNotifEnabled=true")
	}
	if resp.Data.Discord.WeeklyNotifEnabled {
		t.Error("expected weeklyNotifEnabled=false")
	}
}
