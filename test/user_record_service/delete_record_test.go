package user_record_service_test

import (
	"bayar-woy-project/models"
	"bayar-woy-project/testutil"
	urs "bayar-woy-project/user_record_service"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestDeleteUncommittedRecord(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db := testutil.SetupTestDB(t)

	user := models.User{Username: "eve", Password: "x"}
	db.Create(&user)

	record := models.Record{
		Title:       "draft record",
		Amount:      10000,
		OwnerID:     user.ID,
		Type:        "expense",
		IsCommitted: false,
	}
	db.Create(&record)

	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)
	c.Request = httptest.NewRequest(http.MethodDelete, "/user/record/"+record.ID, nil)
	c.Params = gin.Params{{Key: "id", Value: record.ID}}
	c.Set("userID", user.ID)

	urs.DeleteRecord(c)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", rec.Code, rec.Body.String())
	}

	var count int64
	db.Model(&models.Record{}).Where("id = ?", record.ID).Count(&count)
	if count != 0 {
		t.Error("expected record to be deleted")
	}
}

func TestDeleteCommittedRecordReturns409(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db := testutil.SetupTestDB(t)

	user := models.User{Username: "frank", Password: "x"}
	db.Create(&user)

	record := models.Record{
		Title:       "committed record",
		Amount:      5000,
		OwnerID:     user.ID,
		Type:        "expense",
		IsCommitted: true,
	}
	db.Create(&record)

	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)
	c.Request = httptest.NewRequest(http.MethodDelete, "/user/record/"+record.ID, nil)
	c.Params = gin.Params{{Key: "id", Value: record.ID}}
	c.Set("userID", user.ID)

	urs.DeleteRecord(c)

	if rec.Code != http.StatusConflict {
		t.Errorf("expected 409, got %d", rec.Code)
	}
}

func TestDeleteRecordNotOwnerReturns404(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db := testutil.SetupTestDB(t)

	owner := models.User{Username: "grace", Password: "x"}
	other := models.User{Username: "heidi", Password: "x"}
	db.Create(&owner)
	db.Create(&other)

	record := models.Record{
		Title:       "not mine",
		Amount:      5000,
		OwnerID:     owner.ID,
		Type:        "expense",
		IsCommitted: false,
	}
	db.Create(&record)

	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)
	c.Request = httptest.NewRequest(http.MethodDelete, "/user/record/"+record.ID, nil)
	c.Params = gin.Params{{Key: "id", Value: record.ID}}
	c.Set("userID", other.ID)

	urs.DeleteRecord(c)

	if rec.Code != http.StatusNotFound {
		t.Errorf("expected 404, got %d", rec.Code)
	}
}
