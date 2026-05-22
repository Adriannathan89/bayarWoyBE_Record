package service_test

import (
	"bayar-woy-project/models"
	"bayar-woy-project/service"
	"bayar-woy-project/testutil"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestLoadOwedDebtReturnsDebtsWhereUserIsDebtor(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db := testutil.SetupTestDB(t)

	alice := models.User{Username: "alice_owed", Password: "x"}
	bob := models.User{Username: "bob_owed", Password: "x"}
	db.Create(&alice)
	db.Create(&bob)

	db.Create(&models.Debt{
		OwnerID:     bob.ID,
		DebtorID:    alice.ID,
		Amount:      50000,
		Description: "test debt",
		Status:      "pending",
	})

	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)
	c.Request = httptest.NewRequest(http.MethodGet, "/debt/owed", nil)
	c.Set("userID", alice.ID)

	service.LoadOwedDebt(c)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d — body: %s", rec.Code, rec.Body.String())
	}

	var resp map[string]interface{}
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}
	data, ok := resp["data"].(map[string]interface{})
	if !ok {
		t.Fatalf("expected data object in response, got: %v", resp["data"])
	}
	debts, ok := data["debts"].([]interface{})
	if !ok || len(debts) != 1 {
		t.Fatalf("expected 1 owed debt, got: %v", data["debts"])
	}
}

func TestLoadOwedDebtExcludesDebtsUserOwns(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db := testutil.SetupTestDB(t)

	alice := models.User{Username: "alice_excl", Password: "x"}
	bob := models.User{Username: "bob_excl", Password: "x"}
	db.Create(&alice)
	db.Create(&bob)

	// Alice owns this debt (bob owes alice) — should NOT appear in alice's owed list
	db.Create(&models.Debt{
		OwnerID:     alice.ID,
		DebtorID:    bob.ID,
		Amount:      30000,
		Description: "alice owns this",
		Status:      "pending",
	})

	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)
	c.Request = httptest.NewRequest(http.MethodGet, "/debt/owed", nil)
	c.Set("userID", alice.ID)

	service.LoadOwedDebt(c)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}

	var resp map[string]interface{}
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}
	data, ok := resp["data"].(map[string]interface{})
	if !ok {
		t.Fatalf("expected data object in response, got: %v", resp["data"])
	}
	debts := data["debts"]

	if debts != nil {
		if arr, ok := debts.([]interface{}); ok && len(arr) > 0 {
			t.Fatalf("expected empty debts for owner, got %d items", len(arr))
		}
	}
}

func TestLoadOwedDebtExcludesCompletedDebts(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db := testutil.SetupTestDB(t)

	alice := models.User{Username: "alice_completed", Password: "x"}
	bob := models.User{Username: "bob_completed", Password: "x"}
	db.Create(&alice)
	db.Create(&bob)

	// Alice is debtor but the debt is completed — should NOT appear
	db.Create(&models.Debt{
		OwnerID:     bob.ID,
		DebtorID:    alice.ID,
		Amount:      20000,
		Description: "already paid",
		Status:      "completed",
	})

	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)
	c.Request = httptest.NewRequest(http.MethodGet, "/debt/owed", nil)
	c.Set("userID", alice.ID)

	service.LoadOwedDebt(c)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d — body: %s", rec.Code, rec.Body.String())
	}

	var resp map[string]interface{}
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}
	data, ok := resp["data"].(map[string]interface{})
	if !ok {
		t.Fatalf("expected data object in response, got: %v", resp["data"])
	}
	debts := data["debts"]

	if debts != nil {
		if arr, ok := debts.([]interface{}); ok && len(arr) > 0 {
			t.Fatalf("expected completed debt to be excluded, got %d items", len(arr))
		}
	}
}
