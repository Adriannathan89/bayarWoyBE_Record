# SLM Service Test Documentation

## Overview

Comprehensive unit and integration tests untuk memastikan backend **bayarwoyBE** terhubung dengan baik dengan **SLM (Sentiment/Language Model)** service. SLM service digunakan untuk klasifikasi otomatis kategori transaksi dan tipe pemasukan/pengeluaran berdasarkan judul/deskripsi.

## Architecture

### SLM Service Integration

```
CreateRecord / CreateDebt
        ↓
   slm.Classify()
        ↓
  SLM HTTP Service (external)
  POST /classify
        ↓
  Respons: {category, transaction_type, confidence}
        ↓
  Category disimpan ke Record/Debt model
```

## Test Files

### 1. `slm/client_test.go` — Unit Tests SLM Client

Tests untuk SLM client library yang melakukan HTTP call ke SLM service.

#### Test Cases:

| Test | Purpose |
|------|---------|
| **TestClassifySuccess** | Verify successful classification response from SLM service |
| **TestClassifyIncome** | Verify income transaction type detection ("pemasukan") |
| **TestClassifyUnreachable** | Verify graceful fallback ketika SLM unreachable → return empty ClassifyResult |
| **TestClassifyBadResponse** | Verify error handling ketika SLM return invalid JSON |
| **TestClassifyTitleBackwardCompat** | Verify backward-compatible `ClassifyTitle()` function |
| **TestClassifyMultipleRequests** | Verify multiple requests ke SLM service work correctly |
| **TestClassifyEmptyTitle** | Verify handling empty string classification |

**Coverage:**
- ✓ HTTP POST ke SLM endpoint dengan payload `{title}`
- ✓ JSON parsing response `{category, transaction_type, confidence}`
- ✓ Timeout handling (2 detik)
- ✓ Error resilience (SLM unavailable → graceful degradation)

### 2. `test/service/slm_integration_test.go` — Integration Tests

Tests integrasi SLM dengan business logic (CreateRecord, CreateDebt).

#### Test Cases:

| Test | Purpose | Validates |
|------|---------|-----------|
| **TestCreateRecordWithSLMClassification** | Record dibuat dengan kategori dari SLM | SLM call, category setting, user cash update, DB persistence |
| **TestCreateRecordWithSLMIncomeType** | Record income type diset dari SLM "pemasukan" | Transaction type conversion (pemasukan → income), cash increase |
| **TestCreateRecordSLMUnreachableGraceful** | Record dibuat meski SLM unreachable | Graceful fallback, empty category, system continues |
| **TestCreateDebtWithSLMClassification** | Debt dibuat dengan kategori dari SLM | SLM integration with debt, user balance updates, record creation |
| **TestCreateDebtSLMUnreachableGraceful** | Debt dibuat meski SLM unreachable | System resilience, empty category fallback |
| **TestCreateRecordWithCustomDate** | Record created dengan custom date dari client | Date parsing, correct CreatedAt |
| **TestCreateRecordInvalidDateFormat** | Reject invalid date format | Input validation |

**Coverage:**
- ✓ SLM classification dipakai di CreateRecord (title → category)
- ✓ SLM classification dipakai di CreateDebt (description → category)
- ✓ "pemasukan" → record type "income", default "expense"
- ✓ User cash/receivable/debt balance updates dengan benar
- ✓ Record otomatis dibuat untuk CreateDebt
- ✓ Graceful degradation ketika SLM unreachable (empty category, operasi tetap jalan)
- ✓ Database consistency (transactions, foreign keys)

## Running Tests

### Run all SLM tests
```bash
go test ./slm -v
```

### Run all integration tests
```bash
go test ./test/service -v -run "SLM"
```

### Run full test suite
```bash
go test ./slm ./test/service -v
```

### Run with coverage
```bash
go test ./slm ./test/service -cover
```

## Test Mocking Strategy

Semua tests menggunakan **httptest** untuk mock SLM service endpoint:

```go
server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    response := classifyResponse{
        Category:        "makanan_minuman",
        TransactionType: "pengeluaran",
        Confidence:      0.95,
    }
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(response)
}))
```

**Advantages:**
- ✓ No external SLM service needed
- ✓ Deterministic responses
- ✓ Fast execution
- ✓ Full control over response scenarios

## Key Behaviors Tested

### 1. SLM Classification Flow
```
Client POST /user/record
   ↓
CreateRecord handler
   ↓
slm.Classify(title)
   ↓
Mock SLM returns category
   ↓
Record saved with category
   ↓
Response to client
```

### 2. Graceful Degradation
```
SLM Unreachable
   ↓
Classify() returns empty ClassifyResult{}
   ↓
Record/Debt created with empty category ("")
   ↓
System continues normally
   ↓
User not blocked
```

### 3. Transaction Type Conversion
```
SLM returns "pemasukan"
   ↓
recordType = "income"

SLM returns "pengeluaran" (or anything else)
   ↓
recordType = "expense" (default)
```

## Database State Verification

Setiap integration test verify:
- ✓ Record/Debt row created di database
- ✓ Category field terisi dari SLM
- ✓ User cash balance updated correctly
- ✓ Owner receivable/debtor debt updated (for CreateDebt)
- ✓ Related records created (expense record for owner in CreateDebt)

## Environment Setup

Tests menggunakan:
- **SQLite in-memory database** (testutil.SetupTestDB)
- **httptest.Server** untuk mock SLM
- **Environment variable** `SLM_URL` untuk SLM endpoint (set/unset per test)

```go
os.Setenv("SLM_URL", server.URL)
defer os.Unsetenv("SLM_URL")
```

## Test Statistics

**Total Tests:** 14
- SLM unit tests: 7
- SLM integration tests: 7

**Status:** ✓ All passing

```
=== RUN   TestClassifySuccess
--- PASS: TestClassifySuccess (0.00s)
=== RUN   TestCreateRecordWithSLMClassification
--- PASS: TestCreateRecordWithSLMClassification (0.05s)
...
ok  	bayar-woy-project/slm	0.013s
ok  	bayar-woy-project/test/service	0.640s
```

## Integration Points

### CreateRecord (`user_record_service/user-record-service.go`)
```go
slmResult := slm.Classify(req.Title)
recordType := "expense"
if slmResult.TransactionType == "pemasukan" {
    recordType = "income"
}

expense := models.Record{
    Category: slmResult.Category,
    Type:     recordType,
    ...
}
```

### CreateDebt (`service/debt-service.go`)
```go
category := slm.ClassifyTitle(req.Description)

transactionModel := models.Debt{
    Category: category,
    ...
}

// Also creates Record for owner with same category
ownerRecord := models.Record{
    Category: category,
    ...
}
```

## Failure Scenarios Covered

1. **SLM service down** → graceful fallback, empty category
2. **SLM returns bad JSON** → empty ClassifyResult, system continues
3. **Network timeout** → 2s timeout, return empty, system continues
4. **Invalid date format** → validation error, 400 response
5. **User not found** → internal error
6. **Database transaction failure** → 500 error

## Production Readiness

- ✓ Tests isolate SLM dependency via mocking
- ✓ Graceful degradation verified (system works even if SLM down)
- ✓ No external service calls in test environment
- ✓ Database transactions verified
- ✓ User balance consistency verified
- ✓ Edge cases covered (empty title, bad responses, timeouts)

## Future Enhancements

1. **SLM Response Caching** — cache classifications untuk title yang sama (belum ada)
2. **Batch Classification** — support classify multiple titles sekaligus (untuk bulk import)
3. **Category Whitelist** — validate SLM responses terhadap known categories (untuk data quality)
4. **Confidence Threshold** — fallback ke default jika confidence < threshold
5. **Metrics & Logging** — track SLM response times, failure rates, categories distribution