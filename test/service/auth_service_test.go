package service_test

import (
	"bayar-woy-project/models"
	svc "bayar-woy-project/service"
	"bayar-woy-project/testutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func TestLoginSuccessCreatesSession(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db := testutil.SetupTestDB(t)

	hash, _ := bcrypt.GenerateFromPassword([]byte("secret"), bcrypt.DefaultCost)
	user := models.User{Username: "alice", Password: string(hash)}
	if err := db.Create(&user).Error; err != nil {
		t.Fatalf("failed to seed user: %v", err)
	}

	c, rec := newJSONContext(http.MethodPost, "/login", gin.H{
		"username": "alice",
		"password": "secret",
	})
	c.Request.RemoteAddr = "127.0.0.1:1234"

	svc.Login(c)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rec.Code)
	}

	var count int64
	if err := db.Model(&models.Session{}).Count(&count).Error; err != nil {
		t.Fatalf("failed counting sessions: %v", err)
	}
	if count != 1 {
		t.Fatalf("expected 1 session, got %d", count)
	}
}

func TestLoginInvalidPassword(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db := testutil.SetupTestDB(t)

	hash, _ := bcrypt.GenerateFromPassword([]byte("secret"), bcrypt.DefaultCost)
	user := models.User{Username: "alice", Password: string(hash)}
	if err := db.Create(&user).Error; err != nil {
		t.Fatalf("failed to seed user: %v", err)
	}

	c, rec := newJSONContext(http.MethodPost, "/login", gin.H{
		"username": "alice",
		"password": "wrong",
	})

	svc.Login(c)

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("expected status 401, got %d", rec.Code)
	}
}

func TestLogoutDeletesSession(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db := testutil.SetupTestDB(t)

	session := models.Session{
		Username:     "bob",
		UserID:       "u1",
		RefreshToken: "rt-1",
		IPAddress:    "127.0.0.1",
		ExpiresAt:    time.Now().Add(1 * time.Hour),
	}
	if err := db.Create(&session).Error; err != nil {
		t.Fatalf("failed to seed session: %v", err)
	}

	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)
	req := httptest.NewRequest(http.MethodPost, "/logout", nil)
	req.AddCookie(&http.Cookie{Name: "refresh_token", Value: "rt-1"})
	c.Request = req

	svc.Logout(c)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rec.Code)
	}

	var count int64
	if err := db.Model(&models.Session{}).Where("refresh_token = ?", "rt-1").Count(&count).Error; err != nil {
		t.Fatalf("failed counting sessions: %v", err)
	}
	if count != 0 {
		t.Fatalf("expected session to be deleted, count=%d", count)
	}
}
