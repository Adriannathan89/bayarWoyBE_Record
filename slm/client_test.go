package slm

import (
	"bayar-woy-project/config"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestClassifySuccess(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Fatalf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/classify" {
			t.Fatalf("expected /classify, got %s", r.URL.Path)
		}

		var payload map[string]string
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			t.Fatalf("failed to decode request: %v", err)
		}
		if title, ok := payload["title"]; !ok || title != "Beli kopi" {
			t.Fatalf("expected title 'Beli kopi', got %s", title)
		}

		response := classifyResponse{
			Category:        "makanan_minuman",
			TransactionType: "pengeluaran",
			Confidence:      0.95,
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	os.Setenv("SLM_URL", server.URL)
	defer os.Unsetenv("SLM_URL")

	result := Classify("Beli kopi")

	if result.Category != "makanan_minuman" {
		t.Fatalf("expected category 'makanan_minuman', got %s", result.Category)
	}
	if result.TransactionType != "pengeluaran" {
		t.Fatalf("expected type 'pengeluaran', got %s", result.TransactionType)
	}
}

func TestClassifyIncome(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response := classifyResponse{
			Category:        "gaji",
			TransactionType: "pemasukan",
			Confidence:      0.98,
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	os.Setenv("SLM_URL", server.URL)
	defer os.Unsetenv("SLM_URL")

	result := Classify("Gaji bulanan")

	if result.TransactionType != "pemasukan" {
		t.Fatalf("expected type 'pemasukan', got %s", result.TransactionType)
	}
	if result.Category != "gaji" {
		t.Fatalf("expected category 'gaji', got %s", result.Category)
	}
}

func TestClassifyUnreachable(t *testing.T) {
	os.Setenv("SLM_URL", "http://invalid-url-that-does-not-exist")
	defer os.Unsetenv("SLM_URL")

	result := Classify("Beli kopi")

	if result.Category != "" {
		t.Fatalf("expected empty category, got %s", result.Category)
	}
	if result.TransactionType != "" {
		t.Fatalf("expected empty type, got %s", result.TransactionType)
	}
}

func TestClassifyBadResponse(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte("invalid json"))
	}))
	defer server.Close()

	os.Setenv("SLM_URL", server.URL)
	defer os.Unsetenv("SLM_URL")

	result := Classify("Beli kopi")

	if result.Category != "" {
		t.Fatalf("expected empty category on bad response, got %s", result.Category)
	}
}

func TestClassifyTitleBackwardCompat(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response := classifyResponse{
			Category:        "hiburan",
			TransactionType: "pengeluaran",
			Confidence:      0.85,
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	os.Setenv("SLM_URL", server.URL)
	defer os.Unsetenv("SLM_URL")

	category := ClassifyTitle("Nonton bioskop")

	if category != "hiburan" {
		t.Fatalf("expected 'hiburan', got %s", category)
	}
}

func TestClassifyMultipleRequests(t *testing.T) {
	requestCount := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestCount++
		var payload map[string]string
		json.NewDecoder(r.Body).Decode(&payload)

		response := classifyResponse{
			Category:        "test_" + payload["title"],
			TransactionType: "pengeluaran",
			Confidence:      0.9,
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	os.Setenv("SLM_URL", server.URL)
	defer os.Unsetenv("SLM_URL")

	result1 := Classify("request1")
	result2 := Classify("request2")
	result3 := Classify("request3")

	if requestCount != 3 {
		t.Fatalf("expected 3 requests, got %d", requestCount)
	}
	if result1.Category != "test_request1" || result2.Category != "test_request2" || result3.Category != "test_request3" {
		t.Fatalf("unexpected categories: %s, %s, %s", result1.Category, result2.Category, result3.Category)
	}
}

func TestClassifyReadsURLFromConfigGetEnv(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response := classifyResponse{
			Category:        "gaji",
			TransactionType: "pemasukan",
			Confidence:      0.99,
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	os.Setenv("SLM_URL", server.URL)
	defer os.Unsetenv("SLM_URL")

	// Verify config.GetEnv reads the env var (not bypassed by os.Getenv directly)
	if got := config.GetEnv("SLM_URL"); got != server.URL {
		t.Fatalf("config.GetEnv(SLM_URL) = %q, want %q", got, server.URL)
	}

	// Verify Classify() actually hits the server via config.GetEnv
	result := Classify("gaji bulanan")
	if result.Category != "gaji" {
		t.Fatalf("Classify must read SLM_URL from config.GetEnv: expected 'gaji', got %q", result.Category)
	}
	if result.TransactionType != "pemasukan" {
		t.Fatalf("expected 'pemasukan', got %q", result.TransactionType)
	}
}

func TestClassifyEmptyTitle(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response := classifyResponse{
			Category:        "unknown",
			TransactionType: "pengeluaran",
			Confidence:      0.0,
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	os.Setenv("SLM_URL", server.URL)
	defer os.Unsetenv("SLM_URL")

	result := Classify("")

	if result.Category != "unknown" {
		t.Fatalf("expected 'unknown', got %s", result.Category)
	}
}