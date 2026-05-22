package service_test

import (
	"bayar-woy-project/models"
	"bayar-woy-project/service"
	"bayar-woy-project/testutil"
	"bytes"
	"encoding/json"
	"fmt"
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

func TestCreateDebtCreatesExpenseRecordForOwner(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db := testutil.SetupTestDB(t)

	alice := models.User{Username: "alice_create", Password: "x", Cash: 100000}
	bob := models.User{Username: "bob_create", Password: "x"}
	db.Create(&alice)
	db.Create(&bob)

	body := bytes.NewBufferString(fmt.Sprintf(
		`{"amount":50000,"description":"makan","debtorId":"%s"}`, bob.ID,
	))
	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)
	c.Request = httptest.NewRequest(http.MethodPost, "/debt/create", body)
	c.Request.Header.Set("Content-Type", "application/json")
	c.Set("userID", alice.ID)

	service.CreateDebt(c)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d — body: %s", rec.Code, rec.Body.String())
	}

	var records []models.Record
	db.Where("owner_id = ? AND type = ?", alice.ID, "expense").Find(&records)
	if len(records) != 1 {
		t.Fatalf("expected 1 expense record for owner, got %d", len(records))
	}
	if records[0].Amount != 50000 {
		t.Fatalf("expected amount 50000, got %f", records[0].Amount)
	}
	if records[0].Title != "Hutang ke bob_create" {
		t.Fatalf("expected title 'Hutang ke bob_create', got %q", records[0].Title)
	}

	var updatedAlice models.User
	db.First(&updatedAlice, "id = ?", alice.ID)
	if updatedAlice.Cash != 50000 {
		t.Fatalf("expected alice.Cash=50000 (100000-50000), got %f", updatedAlice.Cash)
	}
}

func TestFinishDebtCreatesRecordsAndUpdatesBothBalances(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db := testutil.SetupTestDB(t)

	alice := models.User{Username: "alice_finish", Password: "x", Cash: 0, Receivable: 50000}
	bob := models.User{Username: "bob_finish", Password: "x", Cash: 100000, Debt: 50000}
	db.Create(&alice)
	db.Create(&bob)

	debt := models.Debt{
		OwnerID:     alice.ID,
		DebtorID:    bob.ID,
		Amount:      50000,
		Description: "test",
		Status:      "pending",
	}
	db.Create(&debt)

	body := bytes.NewBufferString(fmt.Sprintf(`{"debtId":"%s"}`, debt.ID))
	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)
	c.Request = httptest.NewRequest(http.MethodPut, "/debt/finish", body)
	c.Request.Header.Set("Content-Type", "application/json")
	c.Set("userID", bob.ID)

	service.FinishDebt(c)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d — body: %s", rec.Code, rec.Body.String())
	}

	var updatedAlice, updatedBob models.User
	db.First(&updatedAlice, "id = ?", alice.ID)
	db.First(&updatedBob, "id = ?", bob.ID)

	if updatedAlice.Cash != 50000 {
		t.Fatalf("expected alice.Cash=50000, got %f", updatedAlice.Cash)
	}
	if updatedAlice.Receivable != 0 {
		t.Fatalf("expected alice.Receivable=0, got %f", updatedAlice.Receivable)
	}
	if updatedBob.Cash != 50000 {
		t.Fatalf("expected bob.Cash=50000 (100000-50000), got %f", updatedBob.Cash)
	}
	if updatedBob.Debt != 0 {
		t.Fatalf("expected bob.Debt=0, got %f", updatedBob.Debt)
	}

	var ownerRecords, debtorRecords []models.Record
	db.Where("owner_id = ? AND type = ?", alice.ID, "income").Find(&ownerRecords)
	db.Where("owner_id = ? AND type = ?", bob.ID, "expense").Find(&debtorRecords)
	if len(ownerRecords) != 1 {
		t.Fatalf("expected 1 income record for owner, got %d", len(ownerRecords))
	}
	if len(debtorRecords) != 1 {
		t.Fatalf("expected 1 expense record for debtor, got %d", len(debtorRecords))
	}

	if ownerRecords[0].Amount != 50000 {
		t.Fatalf("expected owner record amount 50000, got %f", ownerRecords[0].Amount)
	}
	if ownerRecords[0].Title != "Piutang lunas dari bob_finish" {
		t.Fatalf("expected owner record title 'Piutang lunas dari bob_finish', got %q", ownerRecords[0].Title)
	}
	if debtorRecords[0].Amount != 50000 {
		t.Fatalf("expected debtor record amount 50000, got %f", debtorRecords[0].Amount)
	}
	if debtorRecords[0].Title != "Bayar hutang ke alice_finish" {
		t.Fatalf("expected debtor record title 'Bayar hutang ke alice_finish', got %q", debtorRecords[0].Title)
	}
}

func TestFinishDebtRejectsNonDebtor(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db := testutil.SetupTestDB(t)

	alice := models.User{Username: "alice_auth", Password: "x"}
	bob := models.User{Username: "bob_auth", Password: "x"}
	charlie := models.User{Username: "charlie_auth", Password: "x"}
	db.Create(&alice)
	db.Create(&bob)
	db.Create(&charlie)

	debt := models.Debt{
		OwnerID:  alice.ID,
		DebtorID: bob.ID,
		Amount:   10000,
		Status:   "pending",
	}
	db.Create(&debt)

	body := bytes.NewBufferString(fmt.Sprintf(`{"debtId":"%s"}`, debt.ID))
	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)
	c.Request = httptest.NewRequest(http.MethodPut, "/debt/finish", body)
	c.Request.Header.Set("Content-Type", "application/json")
	c.Set("userID", charlie.ID)

	service.FinishDebt(c)

	if rec.Code != http.StatusForbidden {
		t.Fatalf("expected 403, got %d — body: %s", rec.Code, rec.Body.String())
	}
}

func TestFinishDebtRejectsAlreadyCompletedDebt(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db := testutil.SetupTestDB(t)

	alice := models.User{Username: "alice_dbl", Password: "x", Cash: 0, Receivable: 30000}
	bob := models.User{Username: "bob_dbl", Password: "x", Cash: 60000, Debt: 30000}
	db.Create(&alice)
	db.Create(&bob)

	debt := models.Debt{
		OwnerID:  alice.ID,
		DebtorID: bob.ID,
		Amount:   30000,
		Status:   "pending",
	}
	db.Create(&debt)

	callFinish := func() *httptest.ResponseRecorder {
		body := bytes.NewBufferString(fmt.Sprintf(`{"debtId":"%s"}`, debt.ID))
		rec := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(rec)
		c.Request = httptest.NewRequest(http.MethodPut, "/debt/finish", body)
		c.Request.Header.Set("Content-Type", "application/json")
		c.Set("userID", bob.ID)
		service.FinishDebt(c)
		return rec
	}

	first := callFinish()
	if first.Code != http.StatusOK {
		t.Fatalf("first call: expected 200, got %d — %s", first.Code, first.Body.String())
	}

	second := callFinish()
	if second.Code != http.StatusConflict {
		t.Fatalf("second call: expected 409, got %d — %s", second.Code, second.Body.String())
	}

	var bob2 models.User
	db.First(&bob2, "id = ?", bob.ID)
	if bob2.Cash != 30000 {
		t.Fatalf("bob.Cash should be 30000 after single deduction, got %f", bob2.Cash)
	}
}
