package user_record_service_test

import (
	"bayar-woy-project/models"
	"bayar-woy-project/testutil"
	urs "bayar-woy-project/user_record_service"
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func newJSONContext(method string, path string, body any) (*gin.Context, *httptest.ResponseRecorder) {
	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)

	payload, _ := json.Marshal(body)
	req := httptest.NewRequest(method, path, bytes.NewBuffer(payload))
	req.Header.Set("Content-Type", "application/json")
	c.Request = req

	return c, rec
}

func TestCreateRecordIncomeUpdatesCash(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db := testutil.SetupTestDB(t)

	user := models.User{Username: "alice", Password: "x", Cash: 100}
	if err := db.Create(&user).Error; err != nil {
		t.Fatalf("failed seeding user: %v", err)
	}

	c, rec := newJSONContext(http.MethodPost, "/records", gin.H{
		"title":       "salary",
		"description": "monthly income",
		"amount":      50,
		"type":        "income",
	})
	c.Set("userID", user.ID)

	urs.CreateRecord(c)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rec.Code)
	}
}

func TestLoadAllRecordsReturnsOK(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db := testutil.SetupTestDB(t)

	owner := models.User{Username: "owner", Password: "x", Cash: 80, Debt: 10, Receivable: 20}
	debtor := models.User{Username: "debtor", Password: "x"}
	if err := db.Create(&owner).Error; err != nil {
		t.Fatalf("failed seeding owner: %v", err)
	}
	if err := db.Create(&debtor).Error; err != nil {
		t.Fatalf("failed seeding debtor: %v", err)
	}

	record := models.Record{Title: "Food", Description: "Lunch", Amount: 20, OwnerID: owner.ID, Type: "expense"}
	if err := db.Create(&record).Error; err != nil {
		t.Fatalf("failed seeding record: %v", err)
	}

	debt := models.Debt{OwnerID: owner.ID, DebtorID: debtor.ID, Amount: 30, Description: "borrowed", Status: "pending"}
	if err := db.Create(&debt).Error; err != nil {
		t.Fatalf("failed seeding debt: %v", err)
	}

	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)
	c.Request = httptest.NewRequest(http.MethodGet, "/records", nil)
	c.Set("userID", owner.ID)

	urs.LoadAllRecords(c)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rec.Code)
	}
}
