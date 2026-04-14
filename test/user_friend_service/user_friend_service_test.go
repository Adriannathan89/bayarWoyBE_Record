package user_friend_service_test

import (
	"bayar-woy-project/models"
	"bayar-woy-project/testutil"
	ufs "bayar-woy-project/user_friend_service"
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func newJSONContext(method string, path string, body any) (*gin.Context, *httptest.ResponseRecorder) {
	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)

	payload, _ := json.Marshal(body)
	req := httptest.NewRequest(method, path, bytes.NewBuffer(payload))
	req.Header.Set("Content-Type", "application/json")
	c.Request = req

	return c, rec
}

func TestSentFriendRequestSuccess(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db := testutil.SetupTestDB(t)

	sender := models.User{Username: "alice", Password: "x"}
	receiver := models.User{Username: "bob", Password: "x"}
	_ = db.Create(&sender).Error
	_ = db.Create(&receiver).Error

	c, rec := newJSONContext(http.MethodPost, "/friends/request", gin.H{
		"friendUsername": "bob",
	})
	c.Set("username", "alice")

	ufs.SentFriendRequest(c)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rec.Code)
	}
}

func TestGetFriendRequestsReturnsInternalErrorWithCurrentColumnName(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db := testutil.SetupTestDB(t)

	req := models.FriendRequest{SenderID: "alice", ReceiverID: "bob"}
	if err := db.Create(&req).Error; err != nil {
		t.Fatalf("failed seeding request: %v", err)
	}

	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)
	c.Request = httptest.NewRequest(http.MethodGet, "/friends/requests", nil)
	c.Set("userID", "alice")

	ufs.GetFriendRequests(c)

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("expected status 500, got %d", rec.Code)
	}
}

func TestFriendRequestResponseReject(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db := testutil.SetupTestDB(t)

	reqModel := models.FriendRequest{SenderID: "alice", ReceiverID: "bob"}
	if err := db.Create(&reqModel).Error; err != nil {
		t.Fatalf("failed seeding request: %v", err)
	}

	c, rec := newJSONContext(http.MethodPost, "/friends/respond", gin.H{
		"friendshipId": reqModel.ID,
		"action":      "reject",
	})

	ufs.FriendRequestResponse(c)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rec.Code)
	}
}

func TestSearchFriendAndGetAllFriends(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db := testutil.SetupTestDB(t)

	alice := models.User{Username: "alice", Password: "x"}
	bob := models.User{Username: "bobby", Password: "x"}
	if err := db.Create(&alice).Error; err != nil {
		t.Fatalf("failed seeding alice: %v", err)
	}
	if err := db.Create(&bob).Error; err != nil {
		t.Fatalf("failed seeding bob: %v", err)
	}

	friendship := models.Friendship{UserID: alice.ID, FriendID: bob.ID}
	if err := db.Create(&friendship).Error; err != nil {
		t.Fatalf("failed seeding friendship: %v", err)
	}

	searchCtx, searchRec := newJSONContext(http.MethodPost, "/friends/search", gin.H{"name": "bob"})
	searchCtx.Set("userID", alice.ID)
	ufs.SearchFriend(searchCtx)

	if searchRec.Code != http.StatusOK {
		t.Fatalf("expected search status 200, got %d", searchRec.Code)
	}

	listRec := httptest.NewRecorder()
	listCtx, _ := gin.CreateTestContext(listRec)
	listCtx.Request = httptest.NewRequest(http.MethodGet, "/friends", nil)
	listCtx.Set("userID", alice.ID)
	ufs.GetAllFriends(listCtx)

	if listRec.Code != http.StatusOK {
		t.Fatalf("expected list status 200, got %d", listRec.Code)
	}
}
