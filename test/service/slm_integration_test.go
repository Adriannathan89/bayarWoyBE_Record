package service_test

import (
	"bayar-woy-project/dto"
	"bayar-woy-project/models"
	"bayar-woy-project/responses"
	user_record_svc "bayar-woy-project/user_record_service"
	svc "bayar-woy-project/service"
	"bayar-woy-project/testutil"
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func TestCreateRecordWithSLMClassification(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db := testutil.SetupTestDB(t)

	hash, _ := bcrypt.GenerateFromPassword([]byte("secret"), bcrypt.DefaultCost)
	user := models.User{Username: "alice", Password: string(hash), Cash: 1000000}
	if err := db.Create(&user).Error; err != nil {
		t.Fatalf("failed to seed user: %v", err)
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response := map[string]interface{}{
			"category":         "makanan_minuman",
			"transaction_type": "pengeluaran",
			"confidence":       0.95,
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	os.Setenv("SLM_URL", server.URL)
	defer os.Unsetenv("SLM_URL")

	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)
	c.Set("userID", user.ID)

	payload := dto.CreateRecordDto{
		Title:       "Beli kopi",
		Description: "Kopi di kafe",
		Amount:      50000,
		Date:        time.Now().Format("2006-01-02"),
	}
	body, _ := json.Marshal(payload)
	c.Request = httptest.NewRequest(http.MethodPost, "/user/record", bytes.NewBuffer(body))
	c.Request.Header.Set("Content-Type", "application/json")

	user_record_svc.CreateRecord(c)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", rec.Code, rec.Body.String())
	}

	var response responses.APIResponse
	if err := json.NewDecoder(rec.Body).Decode(&response); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	// Verify response contains the category from SLM
	data := response.Data.(map[string]interface{})
	if data["category"] != "makanan_minuman" {
		t.Fatalf("expected category 'makanan_minuman', got %v", data["category"])
	}

	// Verify record was created with correct type (expense)
	if data["type"] != "expense" {
		t.Fatalf("expected type 'expense', got %v", data["type"])
	}

	// Verify record was saved to database
	var record models.Record
	if err := db.Where("owner_id = ? AND title = ?", user.ID, "Beli kopi").First(&record).Error; err != nil {
		t.Fatalf("failed to find record: %v", err)
	}

	if record.Category != "makanan_minuman" {
		t.Fatalf("expected category 'makanan_minuman' in db, got %s", record.Category)
	}

	// Verify user cash was updated (decreased by amount)
	var updatedUser models.User
	db.Where("id = ?", user.ID).First(&updatedUser)
	if updatedUser.Cash != 950000 {
		t.Fatalf("expected cash 950000, got %f", updatedUser.Cash)
	}
}

func TestCreateRecordWithSLMIncomeType(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db := testutil.SetupTestDB(t)

	hash, _ := bcrypt.GenerateFromPassword([]byte("secret"), bcrypt.DefaultCost)
	user := models.User{Username: "bob", Password: string(hash), Cash: 0}
	if err := db.Create(&user).Error; err != nil {
		t.Fatalf("failed to seed user: %v", err)
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response := map[string]interface{}{
			"category":         "gaji",
			"transaction_type": "pemasukan",
			"confidence":       0.98,
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	os.Setenv("SLM_URL", server.URL)
	defer os.Unsetenv("SLM_URL")

	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)
	c.Set("userID", user.ID)

	payload := dto.CreateRecordDto{
		Title:       "Gaji bulanan",
		Description: "Gaji dari perusahaan",
		Amount:      5000000,
		Date:        time.Now().Format("2006-01-02"),
	}
	body, _ := json.Marshal(payload)
	c.Request = httptest.NewRequest(http.MethodPost, "/user/record", bytes.NewBuffer(body))
	c.Request.Header.Set("Content-Type", "application/json")

	user_record_svc.CreateRecord(c)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", rec.Code, rec.Body.String())
	}

	var response responses.APIResponse
	json.NewDecoder(rec.Body).Decode(&response)

	data := response.Data.(map[string]interface{})
	if data["type"] != "income" {
		t.Fatalf("expected type 'income' for pemasukan, got %v", data["type"])
	}

	// Verify user cash was increased
	var updatedUser models.User
	db.Where("id = ?", user.ID).First(&updatedUser)
	if updatedUser.Cash != 5000000 {
		t.Fatalf("expected cash 5000000, got %f", updatedUser.Cash)
	}
}

func TestCreateRecordSLMUnreachableGraceful(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db := testutil.SetupTestDB(t)

	hash, _ := bcrypt.GenerateFromPassword([]byte("secret"), bcrypt.DefaultCost)
	user := models.User{Username: "charlie", Password: string(hash), Cash: 1000}
	if err := db.Create(&user).Error; err != nil {
		t.Fatalf("failed to seed user: %v", err)
	}

	os.Setenv("SLM_URL", "http://invalid-url-does-not-exist")
	defer os.Unsetenv("SLM_URL")

	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)
	c.Set("userID", user.ID)

	payload := dto.CreateRecordDto{
		Title:       "Beli sesuatu",
		Description: "Pembelian tanpa kategori",
		Amount:      100000,
		Date:        time.Now().Format("2006-01-02"),
	}
	body, _ := json.Marshal(payload)
	c.Request = httptest.NewRequest(http.MethodPost, "/user/record", bytes.NewBuffer(body))
	c.Request.Header.Set("Content-Type", "application/json")

	user_record_svc.CreateRecord(c)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected graceful fallback (200), got %d: %s", rec.Code, rec.Body.String())
	}

	var response responses.APIResponse
	json.NewDecoder(rec.Body).Decode(&response)

	data := response.Data.(map[string]interface{})
	if data["category"] != "" {
		t.Fatalf("expected empty category when SLM unavailable, got %v", data["category"])
	}

	// Verify record was still created with empty category
	var record models.Record
	if err := db.Where("owner_id = ? AND title = ?", user.ID, "Beli sesuatu").First(&record).Error; err != nil {
		t.Fatalf("expected record to be created even when SLM unavailable: %v", err)
	}
}

func TestCreateDebtWithSLMClassification(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db := testutil.SetupTestDB(t)

	hash, _ := bcrypt.GenerateFromPassword([]byte("secret"), bcrypt.DefaultCost)
	owner := models.User{Username: "dave", Password: string(hash), Cash: 5000000, Receivable: 0, Debt: 0}
	debtor := models.User{Username: "eve", Password: string(hash), Cash: 1000000, Receivable: 0, Debt: 0}

	if err := db.Create(&owner).Error; err != nil {
		t.Fatalf("failed to seed owner: %v", err)
	}
	if err := db.Create(&debtor).Error; err != nil {
		t.Fatalf("failed to seed debtor: %v", err)
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response := map[string]interface{}{
			"category":         "pinjaman_uang",
			"transaction_type": "pengeluaran",
			"confidence":       0.92,
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	os.Setenv("SLM_URL", server.URL)
	defer os.Unsetenv("SLM_URL")

	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)
	c.Set("userID", owner.ID)

	payload := dto.DebtDTO{
		Amount:      500000,
		Description: "Hutang pembelian barang",
		DebtorID:    debtor.ID,
	}
	body, _ := json.Marshal(payload)
	c.Request = httptest.NewRequest(http.MethodPost, "/debt/create", bytes.NewBuffer(body))
	c.Request.Header.Set("Content-Type", "application/json")

	svc.CreateDebt(c)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", rec.Code, rec.Body.String())
	}

	var response responses.APIResponse
	json.NewDecoder(rec.Body).Decode(&response)

	data := response.Data.(map[string]interface{})
	if data["category"] != "pinjaman_uang" {
		t.Fatalf("expected category 'pinjaman_uang', got %v", data["category"])
	}

	// Verify debt was saved with SLM category
	var debt models.Debt
	if err := db.Where("owner_id = ? AND debtor_id = ?", owner.ID, debtor.ID).First(&debt).Error; err != nil {
		t.Fatalf("failed to find debt: %v", err)
	}

	if debt.Category != "pinjaman_uang" {
		t.Fatalf("expected category 'pinjaman_uang' in debt, got %s", debt.Category)
	}

	// Verify owner receivable increased
	var updatedOwner models.User
	db.Where("id = ?", owner.ID).First(&updatedOwner)
	if updatedOwner.Receivable != 500000 {
		t.Fatalf("expected owner receivable 500000, got %f", updatedOwner.Receivable)
	}

	// Verify debtor debt increased
	var updatedDebtor models.User
	db.Where("id = ?", debtor.ID).First(&updatedDebtor)
	if updatedDebtor.Debt != 500000 {
		t.Fatalf("expected debtor debt 500000, got %f", updatedDebtor.Debt)
	}

	// Verify record was created with SLM category for owner
	var ownerRecord models.Record
	if err := db.Where("owner_id = ?", owner.ID).First(&ownerRecord).Error; err != nil {
		t.Fatalf("failed to find owner record: %v", err)
	}

	if ownerRecord.Category != "pinjaman_uang" {
		t.Fatalf("expected owner record category 'pinjaman_uang', got %s", ownerRecord.Category)
	}
	if ownerRecord.Type != "expense" {
		t.Fatalf("expected owner record type 'expense', got %s", ownerRecord.Type)
	}
}

func TestCreateDebtSLMUnreachableGraceful(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db := testutil.SetupTestDB(t)

	hash, _ := bcrypt.GenerateFromPassword([]byte("secret"), bcrypt.DefaultCost)
	owner := models.User{Username: "frank", Password: string(hash), Cash: 5000000, Receivable: 0, Debt: 0}
	debtor := models.User{Username: "grace", Password: string(hash), Cash: 1000000, Receivable: 0, Debt: 0}

	if err := db.Create(&owner).Error; err != nil {
		t.Fatalf("failed to seed owner: %v", err)
	}
	if err := db.Create(&debtor).Error; err != nil {
		t.Fatalf("failed to seed debtor: %v", err)
	}

	os.Setenv("SLM_URL", "http://invalid-slm-url")
	defer os.Unsetenv("SLM_URL")

	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)
	c.Set("userID", owner.ID)

	payload := dto.DebtDTO{
		Amount:      300000,
		Description: "Hutang tanpa kategori",
		DebtorID:    debtor.ID,
	}
	body, _ := json.Marshal(payload)
	c.Request = httptest.NewRequest(http.MethodPost, "/debt/create", bytes.NewBuffer(body))
	c.Request.Header.Set("Content-Type", "application/json")

	svc.CreateDebt(c)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected graceful fallback (200), got %d: %s", rec.Code, rec.Body.String())
	}

	// Verify debt was still created with empty category
	var debt models.Debt
	if err := db.Where("owner_id = ? AND debtor_id = ?", owner.ID, debtor.ID).First(&debt).Error; err != nil {
		t.Fatalf("expected debt to be created even when SLM unavailable: %v", err)
	}

	if debt.Category != "" {
		t.Fatalf("expected empty category when SLM unavailable, got %s", debt.Category)
	}
}

func TestCreateRecordWithCustomDate(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db := testutil.SetupTestDB(t)

	hash, _ := bcrypt.GenerateFromPassword([]byte("secret"), bcrypt.DefaultCost)
	user := models.User{Username: "henry", Password: string(hash), Cash: 2000}
	if err := db.Create(&user).Error; err != nil {
		t.Fatalf("failed to seed user: %v", err)
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response := map[string]interface{}{
			"category":         "test",
			"transaction_type": "pengeluaran",
			"confidence":       0.8,
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	os.Setenv("SLM_URL", server.URL)
	defer os.Unsetenv("SLM_URL")

	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)
	c.Set("userID", user.ID)

	customDate := "2026-05-20"
	payload := dto.CreateRecordDto{
		Title:       "Beli barang",
		Description: "Pembelian kemarin",
		Amount:      100,
		Date:        customDate,
	}
	body, _ := json.Marshal(payload)
	c.Request = httptest.NewRequest(http.MethodPost, "/user/record", bytes.NewBuffer(body))
	c.Request.Header.Set("Content-Type", "application/json")

	user_record_svc.CreateRecord(c)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", rec.Code, rec.Body.String())
	}

	var record models.Record
	if err := db.Where("owner_id = ? AND title = ?", user.ID, "Beli barang").First(&record).Error; err != nil {
		t.Fatalf("failed to find record: %v", err)
	}

	expectedDate, _ := time.Parse("2006-01-02", customDate)
	if record.CreatedAt.Format("2006-01-02") != expectedDate.Format("2006-01-02") {
		t.Fatalf("expected date %s, got %s", customDate, record.CreatedAt.Format("2006-01-02"))
	}
}

func TestCreateRecordInvalidDateFormat(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db := testutil.SetupTestDB(t)

	hash, _ := bcrypt.GenerateFromPassword([]byte("secret"), bcrypt.DefaultCost)
	user := models.User{Username: "iris", Password: string(hash), Cash: 1000}
	if err := db.Create(&user).Error; err != nil {
		t.Fatalf("failed to seed user: %v", err)
	}

	os.Setenv("SLM_URL", "http://localhost:8000")
	defer os.Unsetenv("SLM_URL")

	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)
	c.Set("userID", user.ID)

	payload := dto.CreateRecordDto{
		Title:       "Test",
		Description: "Test",
		Amount:      100,
		Date:        "invalid-date",
	}
	body, _ := json.Marshal(payload)
	c.Request = httptest.NewRequest(http.MethodPost, "/user/record", bytes.NewBuffer(body))
	c.Request.Header.Set("Content-Type", "application/json")

	user_record_svc.CreateRecord(c)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 for invalid date, got %d", rec.Code)
	}
}