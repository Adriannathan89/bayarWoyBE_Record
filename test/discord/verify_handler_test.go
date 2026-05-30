package discord_test

import (
	"bayar-woy-project/discord"
	"bayar-woy-project/models"
	"bayar-woy-project/testutil"
	"testing"
	"time"
)

func TestProcessVerifyCodeSuccess(t *testing.T) {
	db := testutil.SetupTestDB(t)

	user := models.User{Username: "alice", Password: "x"}
	db.Create(&user)

	verification := models.DiscordVerification{
		UserID:    user.ID,
		Code:      "123456",
		ExpiresAt: time.Now().Add(10 * time.Minute),
	}
	db.Create(&verification)

	resultMsg, err := discord.ProcessVerifyCode("123456", "999000111", "alicediscord")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if resultMsg == "" {
		t.Fatal("expected non-empty success message")
	}

	var updated models.User
	db.First(&updated, "id = ?", user.ID)
	if updated.DiscordID == nil || *updated.DiscordID != "999000111" {
		t.Errorf("expected DiscordID set to '999000111', got %v", updated.DiscordID)
	}
	if updated.DiscordUsername == nil || *updated.DiscordUsername != "alicediscord" {
		t.Errorf("expected DiscordUsername set to 'alicediscord', got %v", updated.DiscordUsername)
	}
	if !updated.IsValidated {
		t.Error("expected IsValidated to be true")
	}

	var remaining models.DiscordVerification
	res := db.First(&remaining, "code = ?", "123456")
	if res.Error == nil {
		t.Error("expected verification row to be deleted after success")
	}
}

func TestProcessVerifyCodeExpired(t *testing.T) {
	db := testutil.SetupTestDB(t)

	user := models.User{Username: "bob", Password: "x"}
	db.Create(&user)

	verification := models.DiscordVerification{
		UserID:    user.ID,
		Code:      "654321",
		ExpiresAt: time.Now().Add(-1 * time.Minute), // expired
	}
	db.Create(&verification)

	_, err := discord.ProcessVerifyCode("654321", "999000111", "bobdiscord")
	if err == nil {
		t.Fatal("expected error for expired code")
	}
	if err != discord.ErrCodeExpired {
		t.Errorf("expected ErrCodeExpired, got %v", err)
	}
}

func TestProcessVerifyCodeNotFound(t *testing.T) {
	testutil.SetupTestDB(t)

	_, err := discord.ProcessVerifyCode("000000", "999000111", "discord")
	if err != discord.ErrCodeNotFound {
		t.Errorf("expected ErrCodeNotFound, got %v", err)
	}
}

func TestProcessVerifyCodeDiscordIDAlreadyLinked(t *testing.T) {
	db := testutil.SetupTestDB(t)

	otherDiscordID := "999000111"
	existingUser := models.User{Username: "carol", Password: "x", DiscordID: &otherDiscordID, IsValidated: true}
	db.Create(&existingUser)

	newUser := models.User{Username: "dan", Password: "x"}
	db.Create(&newUser)

	verification := models.DiscordVerification{
		UserID:    newUser.ID,
		Code:      "111222",
		ExpiresAt: time.Now().Add(10 * time.Minute),
	}
	db.Create(&verification)

	_, err := discord.ProcessVerifyCode("111222", "999000111", "discorduser")
	if err != discord.ErrDiscordAlreadyLinked {
		t.Errorf("expected ErrDiscordAlreadyLinked, got %v", err)
	}
}

func TestProcessVerifyCodeUserAlreadyLinked(t *testing.T) {
	db := testutil.SetupTestDB(t)

	existingDiscordID := "888777666"
	user := models.User{Username: "eve", Password: "x", DiscordID: &existingDiscordID, IsValidated: true}
	db.Create(&user)

	verification := models.DiscordVerification{
		UserID:    user.ID,
		Code:      "333444",
		ExpiresAt: time.Now().Add(10 * time.Minute),
	}
	db.Create(&verification)

	_, err := discord.ProcessVerifyCode("333444", "999000111", "discorduser")
	if err != discord.ErrUserAlreadyLinked {
		t.Errorf("expected ErrUserAlreadyLinked, got %v", err)
	}
}
