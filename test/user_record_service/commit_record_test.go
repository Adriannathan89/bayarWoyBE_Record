package user_record_service_test

import (
	"bayar-woy-project/models"
	"bayar-woy-project/testutil"
	urs "bayar-woy-project/user_record_service"
	"net/http"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestCommitRecordUpdatesCashAndMarksCommitted(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db := testutil.SetupTestDB(t)

	user := models.User{Username: "bob", Password: "x", Cash: 0}
	db.Create(&user)

	primary := models.Category{Name: "makanan", Type: "primary"}
	secondary := models.Category{Name: "makanan", Type: "secondary"}
	db.Create(&primary)
	db.Create(&secondary)

	record := models.Record{
		Title:       "nasi goreng",
		Amount:      25000,
		OwnerID:     user.ID,
		Type:        "expense",
		IsCommitted: false,
		Categories:  []models.Category{primary, secondary},
	}
	db.Create(&record)

	c, rec := newJSONContext(http.MethodPut, "/user/record/commit", gin.H{
		"recordId": record.ID,
	})
	c.Set("userID", user.ID)

	urs.CommitRecord(c)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", rec.Code, rec.Body.String())
	}

	var updated models.Record
	db.First(&updated, "id = ?", record.ID)
	if !updated.IsCommitted {
		t.Error("expected IsCommitted to be true after commit")
	}

	var updatedUser models.User
	db.First(&updatedUser, "id = ?", user.ID)
	if updatedUser.Cash != -25000 {
		t.Errorf("expected cash -25000 (expense deducted), got %.2f", updatedUser.Cash)
	}
}

func TestCommitRecordWithCorrectedCategoryUpdatesRecord(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db := testutil.SetupTestDB(t)

	user := models.User{Username: "carol", Password: "x", Cash: 0}
	db.Create(&user)

	oldPrimary := models.Category{Name: "makanan", Type: "primary"}
	oldSecondary := models.Category{Name: "makanan", Type: "secondary"}
	newPrimary := models.Category{Name: "transport", Type: "primary"}
	newSecondary := models.Category{Name: "transport", Type: "secondary"}
	db.Create(&oldPrimary)
	db.Create(&oldSecondary)
	db.Create(&newPrimary)
	db.Create(&newSecondary)

	record := models.Record{
		Title:       "grab ke kantor",
		Amount:      15000,
		OwnerID:     user.ID,
		Type:        "expense",
		IsCommitted: false,
		Categories:  []models.Category{oldPrimary, oldSecondary},
	}
	db.Create(&record)

	c, rec := newJSONContext(http.MethodPut, "/user/record/commit", gin.H{
		"recordId": record.ID,
		"category": "transport",
	})
	c.Set("userID", user.ID)

	urs.CommitRecord(c)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", rec.Code, rec.Body.String())
	}

	var updated models.Record
	db.Preload("Categories").First(&updated, "id = ?", record.ID)
	if !updated.IsCommitted {
		t.Error("expected IsCommitted true")
	}
	found := false
	for _, cat := range updated.Categories {
		if cat.Name == "transport" && cat.Type == "primary" {
			found = true
		}
	}
	if !found {
		t.Error("expected category to be updated to 'transport'")
	}
}

func TestCommitRecordAlreadyCommittedReturns409(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db := testutil.SetupTestDB(t)

	user := models.User{Username: "dave", Password: "x", Cash: 0}
	db.Create(&user)

	record := models.Record{
		Title:       "already committed",
		Amount:      5000,
		OwnerID:     user.ID,
		Type:        "expense",
		IsCommitted: true,
	}
	db.Create(&record)

	c, rec := newJSONContext(http.MethodPut, "/user/record/commit", gin.H{
		"recordId": record.ID,
	})
	c.Set("userID", user.ID)

	urs.CommitRecord(c)

	if rec.Code != http.StatusConflict {
		t.Errorf("expected 409, got %d", rec.Code)
	}
}

func TestCommitRecordNotOwnerReturns404(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db := testutil.SetupTestDB(t)

	owner := models.User{Username: "owner2", Password: "x"}
	other := models.User{Username: "other2", Password: "x"}
	db.Create(&owner)
	db.Create(&other)

	record := models.Record{
		Title:       "private record",
		Amount:      5000,
		OwnerID:     owner.ID,
		Type:        "expense",
		IsCommitted: false,
	}
	db.Create(&record)

	c, rec := newJSONContext(http.MethodPut, "/user/record/commit", gin.H{
		"recordId": record.ID,
	})
	c.Set("userID", other.ID)

	urs.CommitRecord(c)

	if rec.Code != http.StatusNotFound {
		t.Errorf("expected 404, got %d", rec.Code)
	}
}
