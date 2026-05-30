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

func TestGetDiscordStatusNotVerified(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db := testutil.SetupTestDB(t)

	user := models.User{Username: "alice", Password: "x"}
	db.Create(&user)

	c, rec := newJSONContext(http.MethodGet, "/user/discord/status", nil)
	c.Set("userID", user.ID)

	ups.GetDiscordStatus(c)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}

	var resp struct {
		Data responses.DiscordStatusResponse `json:"data"`
	}
	json.Unmarshal(rec.Body.Bytes(), &resp)

	if resp.Data.Verified {
		t.Error("expected verified=false")
	}
	if resp.Data.DiscordUsername != "" {
		t.Errorf("expected empty discordUsername, got %q", resp.Data.DiscordUsername)
	}
}

func TestGetDiscordStatusVerified(t *testing.T) {
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

	c, rec := newJSONContext(http.MethodGet, "/user/discord/status", nil)
	c.Set("userID", user.ID)

	ups.GetDiscordStatus(c)

	var resp struct {
		Data responses.DiscordStatusResponse `json:"data"`
	}
	json.Unmarshal(rec.Body.Bytes(), &resp)

	if !resp.Data.Verified {
		t.Error("expected verified=true")
	}
	if resp.Data.DiscordUsername != "alicediscord" {
		t.Errorf("expected discordUsername='alicediscord', got %q", resp.Data.DiscordUsername)
	}
}
