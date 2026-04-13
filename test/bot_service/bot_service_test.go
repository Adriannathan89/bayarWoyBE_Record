package botservice_test

import (
	botservice "bayar-woy-project/bot-service"
	"bayar-woy-project/models"
	"bayar-woy-project/testutil"
	"strings"
	"testing"
)

func TestRegisterUserToBotUserNotFound(t *testing.T) {
	testutil.SetupTestDB(t)

	res := botservice.RegisterUserToBot("missing-user", "discord-1")
	if res != "User not found" {
		t.Fatalf("expected user not found message, got %s", res)
	}
}

func TestRegisterUserToBotSuccess(t *testing.T) {
	db := testutil.SetupTestDB(t)

	user := models.User{Username: "alice", Password: "x"}
	if err := db.Create(&user).Error; err != nil {
		t.Fatalf("failed seeding user: %v", err)
	}

	res := botservice.RegisterUserToBot("alice", "discord-1")
	if !strings.Contains(res, "otp code:") {
		t.Fatalf("expected otp response, got %s", res)
	}
}

func TestGetFriendsListUserNotFound(t *testing.T) {
	testutil.SetupTestDB(t)

	res := botservice.GetFriendsList("missing-discord")
	if res != "User not found" {
		t.Fatalf("expected user not found, got %s", res)
	}
}
