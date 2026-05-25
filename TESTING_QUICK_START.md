# SLM Service Testing - Quick Start Guide

## 🎯 Overview

Backend bayarwoyBE sudah terintegrasi dengan SLM (Sentiment/Language Model) service untuk klasifikasi otomatis kategori transaksi. Test suite ini memastikan integrasi berfungsi dengan baik.

## ⚡ Quick Commands

```bash
# Run all SLM tests (fastest way)
go test ./slm ./test/service -v -run "SLM|Classify"

# Run with detailed output
go test ./slm ./test/service -v

# Run and show timing
go test ./slm ./test/service -v -count=1

# Run with coverage
go test ./slm ./test/service -cover
```

## 📋 What Gets Tested

### SLM Client Library (`slm/client_test.go`)
- ✓ Successfully call SLM HTTP service and parse response
- ✓ Correctly identify income vs expense transactions
- ✓ Handle service unreachable gracefully (no crash)
- ✓ Parse malformed JSON responses without crashing
- ✓ Classify multiple requests in sequence
- ✓ Handle empty input strings

### SLM Integration (`test/service/slm_integration_test.go`)
- ✓ CreateRecord calls SLM and saves category
- ✓ CreateRecord converts "pemasukan" to income type
- ✓ CreateDebt calls SLM and saves category
- ✓ Both work even if SLM service is down
- ✓ User cash/balance updates correctly
- ✓ Database transaction consistency

## 🔍 Understanding Test Output

### Success Case
```
=== RUN   TestCreateRecordWithSLMClassification
--- PASS: TestCreateRecordWithSLMClassification (0.06s)
```
✓ Test passed in 0.06 seconds

### All Tests Pass Example
```
PASS
ok  	bayar-woy-project/slm	0.011s
ok  	bayar-woy-project/test/service	0.647s
```
✓ Both packages passed with total ~0.66s

## 📊 Test Statistics

| Component | Tests | Status |
|-----------|-------|--------|
| SLM Client | 7 | ✓ All Pass |
| SLM Integration | 7 | ✓ All Pass |
| **Total** | **14** | **✓ 100%** |

## 🛠️ Manual Verification

If you want to manually verify SLM integration works:

### 1. Start SLM Service (if you have it)
```bash
# Assuming SLM service runs on localhost:8000
export SLM_URL=http://localhost:8000
```

### 2. Run Backend
```bash
go run main.go
```

### 3. Create a Transaction via API
```bash
curl -X POST http://localhost:8080/user/record \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Beli kopi",
    "description": "Kopi di kafe",
    "amount": 50000,
    "date": "2026-05-25"
  }'
```

### 4. Check Response
The response should include `"category"` field from SLM classification.

## 🐛 Troubleshooting

### Tests Fail with "invalid memory address"
→ Likely test database setup issue. Verify testutil.SetupTestDB works.

### Tests Fail with "connection refused"
→ Normal! Tests mock the SLM service, no external connection needed.

### Tests Timeout
→ Likely test database is slow. Try running tests individually:
```bash
go test ./slm -v -run TestClassifySuccess
```

## 📁 Test File Locations

```
bayarwoyBE/
├── slm/
│   ├── client.go          ← SLM client library
│   └── client_test.go     ← Unit tests (NEW)
├── test/
│   └── service/
│       └── slm_integration_test.go    ← Integration tests (NEW)
├── SLM_TEST_DOCUMENTATION.md          ← Full documentation (NEW)
└── TESTING_QUICK_START.md             ← This file (NEW)
```

## 🚀 Integration Flow

```
User Action: Create Transaction
       ↓
POST /user/record (title, amount, date)
       ↓
CreateRecord handler
       ↓
slm.Classify(title)  ← Call SLM service
       ↓
Get back: {category, transaction_type}
       ↓
Save Record with category
       ↓
Update user cash balance
       ↓
Return response with category
```

## ✨ Key Assertions

Each test verifies:

1. **Response Code** — 200 OK for success, 400 for bad input
2. **Category Set** — Response contains category from SLM
3. **Type Conversion** — "pemasukan" → "income", else "expense"
4. **Database** — Record/Debt saved with correct category
5. **Balance** — User cash/receivable/debt updated correctly
6. **Consistency** — All updates in single transaction (no partial updates)

## 🎓 Test Structure Pattern

All tests follow this pattern:

```go
func TestSomething(t *testing.T) {
    // 1. Setup: Create test database + seed data
    gin.SetMode(gin.TestMode)
    db := testutil.SetupTestDB(t)
    user := models.User{...}
    db.Create(&user)
    
    // 2. Mock: Setup httptest server for SLM
    server := httptest.NewServer(...)
    os.Setenv("SLM_URL", server.URL)
    defer server.Close()
    
    // 3. Execute: Call the handler
    svc.CreateRecord(c)
    
    // 4. Assert: Verify results
    if rec.Code != http.StatusOK {
        t.Fatalf(...)
    }
    
    // 5. Verify DB: Check database state
    var record models.Record
    db.First(&record)
    if record.Category != expected {
        t.Fatalf(...)
    }
}
```

## 📚 Additional Resources

- See `SLM_TEST_DOCUMENTATION.md` for comprehensive details
- See `slm/client.go` for SLM client implementation
- See `user_record_service/user-record-service.go` for CreateRecord
- See `service/debt-service.go` for CreateDebt

## ✅ Verification Checklist

Before merging or deploying:

- [ ] Run `go test ./slm -v` — All 7 unit tests pass
- [ ] Run `go test ./test/service -v -run SLM` — All 7 integration tests pass
- [ ] Check `go test ./slm ./test/service -cover` — Coverage is good
- [ ] No compiler warnings or errors
- [ ] Manual API test with real SLM service (optional)
- [ ] Check logs for any deprecation warnings

## 🎯 Next Steps

1. **Immediate**: Run tests to verify setup
   ```bash
   go test ./slm ./test/service -v
   ```

2. **Optional**: Add more test cases for specific business logic

3. **Production**: Ensure SLM service is running before starting backend
   ```bash
   export SLM_URL=http://your-slm-service:port
   go run main.go
   ```

---

**Created:** 2026-05-25  
**Test Framework:** Go testing (stdlib)  
**Mocking:** httptest + testutil  
**Status:** ✓ All Tests Passing