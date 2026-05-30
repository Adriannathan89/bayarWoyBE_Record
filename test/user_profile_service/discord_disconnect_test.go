package user_profile_service_test

import (
	"bayar-woy-project/models"
	"bayar-woy-project/testutil"
	ups "bayar-woy-project/user_profile_service"
	"net/http"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestDisconnectDiscordClearsFields(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db := testutil.SetupTestDB(t)

	did := "999000111"
	dun := "alicediscord"
	user := models.User{
		Username:        "alice",
		Password:        "x",
		DiscordID:       &did,
		DiscordUsername: &dun,
		IsValidated:     true,
	}
	db.Create(&user)

	c, rec := newJSONContext(http.MethodDelete, "/user/discord", nil)
	c.Set("userID", user.ID)

	ups.DisconnectDiscord(c)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", rec.Code, rec.Body.String())
	}

	var updated models.User
	db.First(&updated, "id = ?", user.ID)
	if updated.DiscordID != nil {
		t.Errorf("expected DiscordID to be nil, got %v", *updated.DiscordID)
	}
	if updated.DiscordUsername != nil {
		t.Errorf("expected DiscordUsername to be nil, got %v", *updated.DiscordUsername)
	}
	if updated.IsValidated {
		t.Error("expected IsValidated=false after disconnect")
	}
}

func TestDisconnectDiscordAlsoDeletesVerification(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db := testutil.SetupTestDB(t)

	user := models.User{Username: "bob", Password: "x"}
	db.Create(&user)

	db.Create(&models.DiscordVerification{
		UserID: user.ID,
		Code:   "111111",
	})

	c, _ := newJSONContext(http.MethodDelete, "/user/discord", nil)
	c.Set("userID", user.ID)

	ups.DisconnectDiscord(c)

	var count int64
	db.Model(&models.DiscordVerification{}).Where("user_id = ?", user.ID).Count(&count)
	if count != 0 {
		t.Errorf("expected pending verification rows deleted, got %d", count)
	}
}
