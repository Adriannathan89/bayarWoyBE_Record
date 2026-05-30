package user_profile_service_test

import (
	"bayar-woy-project/models"
	"bayar-woy-project/testutil"
	ups "bayar-woy-project/user_profile_service"
	"net/http"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestUpdateNotifSettingsToggleCommitOff(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db := testutil.SetupTestDB(t)

	user := models.User{
		Username: "alice",
		Password: "x",
	}
	db.Create(&user)

	c, rec := newJSONContext(http.MethodPut, "/user/discord/notification", gin.H{
		"type":    "commit",
		"enabled": false,
	})
	c.Set("userID", user.ID)

	ups.UpdateNotifSettings(c)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", rec.Code, rec.Body.String())
	}

	var updated models.User
	db.First(&updated, "id = ?", user.ID)
	if updated.DiscordCommitNotifEnabled {
		t.Error("expected DiscordCommitNotifEnabled=false")
	}
	if !updated.DiscordWeeklyNotifEnabled {
		t.Error("expected DiscordWeeklyNotifEnabled to remain true")
	}
}

func TestUpdateNotifSettingsToggleWeeklyOff(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db := testutil.SetupTestDB(t)

	user := models.User{
		Username: "bob",
		Password: "x",
	}
	db.Create(&user)

	c, _ := newJSONContext(http.MethodPut, "/user/discord/notification", gin.H{
		"type":    "weekly",
		"enabled": false,
	})
	c.Set("userID", user.ID)

	ups.UpdateNotifSettings(c)

	var updated models.User
	db.First(&updated, "id = ?", user.ID)
	if updated.DiscordWeeklyNotifEnabled {
		t.Error("expected DiscordWeeklyNotifEnabled=false")
	}
	if !updated.DiscordCommitNotifEnabled {
		t.Error("expected DiscordCommitNotifEnabled to remain true")
	}
}

func TestUpdateNotifSettingsRejectsInvalidType(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db := testutil.SetupTestDB(t)

	user := models.User{Username: "carol", Password: "x"}
	db.Create(&user)

	c, rec := newJSONContext(http.MethodPut, "/user/discord/notification", gin.H{
		"type":    "invalid",
		"enabled": true,
	})
	c.Set("userID", user.ID)

	ups.UpdateNotifSettings(c)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", rec.Code)
	}
}
