package service_test

import (
	botmodel "bayar-woy-project/bot-model"
	"bayar-woy-project/models"
	svc "bayar-woy-project/service"
	"bayar-woy-project/testutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
)

func TestRegisterSuccess(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db := testutil.SetupTestDB(t)

	c, rec := newJSONContext(http.MethodPost, "/register", gin.H{
		"username": "new-user",
		"password": "secret",
	})

	svc.Register(c)

	if rec.Code != http.StatusCreated {
		t.Fatalf("expected status 201, got %d", rec.Code)
	}

	var count int64
	if err := db.Model(&models.User{}).Where("username = ?", "new-user").Count(&count).Error; err != nil {
		t.Fatalf("failed counting users: %v", err)
	}
	if count != 1 {
		t.Fatalf("expected 1 user, got %d", count)
	}
}

func TestValidateOtpSuccess(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db := testutil.SetupTestDB(t)

	user := models.User{Username: "otp-user", Password: "x"}
	if err := db.Create(&user).Error; err != nil {
		t.Fatalf("failed to seed user: %v", err)
	}

	otp := botmodel.DiscordBotOtp{
		UserID:    user.ID,
		OTP:       "123456",
		ExpiredAt: time.Now().Add(5 * time.Minute),
	}
	if err := db.Create(&otp).Error; err != nil {
		t.Fatalf("failed to seed otp: %v", err)
	}

	c, rec := newJSONContext(http.MethodPost, "/otp", gin.H{
		"userId": user.ID,
		"otp":    "123456",
	})

	svc.ValidateOtp(c)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rec.Code)
	}
}

func TestCreateDebtSuccess(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db := testutil.SetupTestDB(t)

	owner := models.User{Username: "owner", Password: "x", Cash: 100}
	debtor := models.User{Username: "debtor", Password: "x", Cash: 40}
	if err := db.Create(&owner).Error; err != nil {
		t.Fatalf("failed to seed owner: %v", err)
	}
	if err := db.Create(&debtor).Error; err != nil {
		t.Fatalf("failed to seed debtor: %v", err)
	}

	c, rec := newJSONContext(http.MethodPost, "/debts", gin.H{
		"amount":      50,
		"description": "shared dinner",
		"debtorId":    debtor.ID,
	})
	c.Set("userId", owner.ID)

	svc.CreateDebt(c)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rec.Code)
	}
}

func TestFinishDebtSuccess(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db := testutil.SetupTestDB(t)

	owner := models.User{Username: "owner", Password: "x", Receivable: 80, Cash: 10}
	debtor := models.User{Username: "debtor", Password: "x"}
	if err := db.Create(&owner).Error; err != nil {
		t.Fatalf("failed to seed owner: %v", err)
	}
	if err := db.Create(&debtor).Error; err != nil {
		t.Fatalf("failed to seed debtor: %v", err)
	}

	debt := models.Debt{OwnerID: owner.ID, DebtorID: debtor.ID, Amount: 30, Description: "lunch", Status: "pending"}
	if err := db.Create(&debt).Error; err != nil {
		t.Fatalf("failed to seed debt: %v", err)
	}

	c, rec := newJSONContext(http.MethodPost, "/debts/finish", gin.H{
		"debtId": debt.ID,
	})

	svc.FinishDebt(c)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rec.Code)
	}
}

func TestLoadAllDebtReturnsInternalErrorWithCurrentPreloadQuery(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db := testutil.SetupTestDB(t)

	owner := models.User{Username: "owner", Password: "x"}
	debtor := models.User{Username: "debtor", Password: "x"}
	if err := db.Create(&owner).Error; err != nil {
		t.Fatalf("failed to seed owner: %v", err)
	}
	if err := db.Create(&debtor).Error; err != nil {
		t.Fatalf("failed to seed debtor: %v", err)
	}

	debt := models.Debt{OwnerID: owner.ID, DebtorID: debtor.ID, Amount: 25, Description: "coffee", Status: "pending"}
	if err := db.Create(&debt).Error; err != nil {
		t.Fatalf("failed to seed debt: %v", err)
	}

	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)
	c.Request = httptest.NewRequest(http.MethodGet, "/debts", nil)
	c.Set("userID", owner.ID)

	svc.LoadAllDebt(c)

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("expected status 500, got %d", rec.Code)
	}
}
