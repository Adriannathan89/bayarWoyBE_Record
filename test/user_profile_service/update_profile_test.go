package user_profile_service_test

import (
	"bayar-woy-project/models"
	"bayar-woy-project/testutil"
	ups "bayar-woy-project/user_profile_service"
	"net/http"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestUpdateProfileChangesUsername(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db := testutil.SetupTestDB(t)

	user := models.User{Username: "oldname", Password: "x"}
	db.Create(&user)

	c, rec := newJSONContext(http.MethodPut, "/user/profile", gin.H{
		"username": "newname",
	})
	c.Set("userID", user.ID)

	ups.UpdateProfile(c)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", rec.Code, rec.Body.String())
	}

	var updated models.User
	db.First(&updated, "id = ?", user.ID)
	if updated.Username != "newname" {
		t.Errorf("expected username 'newname', got %q", updated.Username)
	}
}

func TestUpdateProfileRejectsConflictingUsername(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db := testutil.SetupTestDB(t)

	db.Create(&models.User{Username: "existing", Password: "x"})
	user := models.User{Username: "alice", Password: "x"}
	db.Create(&user)

	c, rec := newJSONContext(http.MethodPut, "/user/profile", gin.H{
		"username": "existing",
	})
	c.Set("userID", user.ID)

	ups.UpdateProfile(c)

	if rec.Code != http.StatusConflict {
		t.Errorf("expected 409, got %d: %s", rec.Code, rec.Body.String())
	}
}

func TestUpdateProfileRejectsEmptyUsername(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db := testutil.SetupTestDB(t)

	user := models.User{Username: "alice", Password: "x"}
	db.Create(&user)

	c, rec := newJSONContext(http.MethodPut, "/user/profile", gin.H{
		"username": "",
	})
	c.Set("userID", user.ID)

	ups.UpdateProfile(c)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", rec.Code)
	}
}
