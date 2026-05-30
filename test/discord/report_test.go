package discord_test

import (
	"bayar-woy-project/discord"
	"bayar-woy-project/models"
	"bayar-woy-project/testutil"
	"testing"
	"time"
)

func TestBuildWeeklyReportEmpty(t *testing.T) {
	db := testutil.SetupTestDB(t)
	user := models.User{Username: "alice", Password: "x"}
	db.Create(&user)

	loc, _ := time.LoadLocation("Asia/Jakarta")
	now := time.Date(2026, 5, 30, 20, 0, 0, 0, loc) // Saturday 20:00 WIB

	report := discord.BuildWeeklyReport(user.ID, now, loc)

	if report.TotalIncome != 0 || report.TotalExpense != 0 {
		t.Errorf("expected zero totals, got income=%f expense=%f", report.TotalIncome, report.TotalExpense)
	}
	if len(report.TopExpenseCategories) != 0 {
		t.Errorf("expected no categories, got %d", len(report.TopExpenseCategories))
	}
}

func TestBuildWeeklyReportAggregations(t *testing.T) {
	db := testutil.SetupTestDB(t)
	user := models.User{Username: "bob", Password: "x"}
	db.Create(&user)

	loc, _ := time.LoadLocation("Asia/Jakarta")
	now := time.Date(2026, 5, 30, 20, 0, 0, 0, loc) // Saturday 20:00 WIB

	// Window is Sunday 24 May 00:00 → Saturday 30 May 19:59:59 WIB
	insideWindow := time.Date(2026, 5, 26, 10, 0, 0, 0, loc) // Tuesday inside

	makanan := models.Category{Name: "makanan", Type: "primary"}
	transport := models.Category{Name: "transport", Type: "primary"}
	gaji := models.Category{Name: "gaji", Type: "primary"}
	db.Create(&makanan)
	db.Create(&transport)
	db.Create(&gaji)

	// Income 5,000,000
	db.Create(&models.Record{Title: "salary", Amount: 5000000, OwnerID: user.ID, Type: "income", IsCommitted: true, CreatedAt: insideWindow, Categories: []models.Category{gaji}})
	// Expense makanan 500,000 (top 1)
	db.Create(&models.Record{Title: "groceries", Amount: 500000, OwnerID: user.ID, Type: "expense", IsCommitted: true, CreatedAt: insideWindow, Categories: []models.Category{makanan}})
	// Expense transport 300,000 (top 2)
	db.Create(&models.Record{Title: "grab", Amount: 300000, OwnerID: user.ID, Type: "expense", IsCommitted: true, CreatedAt: insideWindow, Categories: []models.Category{transport}})
	// Uncommitted — should be ignored
	db.Create(&models.Record{Title: "draft", Amount: 999999, OwnerID: user.ID, Type: "expense", IsCommitted: false, CreatedAt: insideWindow, Categories: []models.Category{makanan}})

	report := discord.BuildWeeklyReport(user.ID, now, loc)

	if report.TotalIncome != 5000000 {
		t.Errorf("expected income 5000000, got %f", report.TotalIncome)
	}
	if report.TotalExpense != 800000 {
		t.Errorf("expected expense 800000, got %f", report.TotalExpense)
	}
	if len(report.TopExpenseCategories) != 2 {
		t.Fatalf("expected 2 categories, got %d", len(report.TopExpenseCategories))
	}
	if report.TopExpenseCategories[0].Name != "makanan" || report.TopExpenseCategories[0].Total != 500000 {
		t.Errorf("expected makanan/500000 first, got %s/%f", report.TopExpenseCategories[0].Name, report.TopExpenseCategories[0].Total)
	}
	if report.TopExpenseCategories[1].Name != "transport" || report.TopExpenseCategories[1].Total != 300000 {
		t.Errorf("expected transport/300000 second, got %s/%f", report.TopExpenseCategories[1].Name, report.TopExpenseCategories[1].Total)
	}
}

func TestBuildWeeklyReportRespectsWindow(t *testing.T) {
	db := testutil.SetupTestDB(t)
	user := models.User{Username: "carol", Password: "x"}
	db.Create(&user)

	loc, _ := time.LoadLocation("Asia/Jakarta")
	now := time.Date(2026, 5, 30, 20, 0, 0, 0, loc)

	makanan := models.Category{Name: "makanan", Type: "primary"}
	db.Create(&makanan)

	beforeWindow := time.Date(2026, 5, 23, 23, 59, 0, 0, loc)  // Last Saturday — outside
	afterWindow := time.Date(2026, 5, 30, 20, 30, 0, 0, loc)    // After 20:00 cutoff — outside
	insideWindow := time.Date(2026, 5, 28, 12, 0, 0, 0, loc)    // Thursday — inside

	db.Create(&models.Record{Title: "old", Amount: 100000, OwnerID: user.ID, Type: "expense", IsCommitted: true, CreatedAt: beforeWindow, Categories: []models.Category{makanan}})
	db.Create(&models.Record{Title: "future", Amount: 200000, OwnerID: user.ID, Type: "expense", IsCommitted: true, CreatedAt: afterWindow, Categories: []models.Category{makanan}})
	db.Create(&models.Record{Title: "in", Amount: 300000, OwnerID: user.ID, Type: "expense", IsCommitted: true, CreatedAt: insideWindow, Categories: []models.Category{makanan}})

	report := discord.BuildWeeklyReport(user.ID, now, loc)

	if report.TotalExpense != 300000 {
		t.Errorf("expected only inside-window 300000, got %f", report.TotalExpense)
	}
}
