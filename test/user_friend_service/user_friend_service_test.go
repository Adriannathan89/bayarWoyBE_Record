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
		"friendId": receiver.ID,
	})
	c.Set("userID", sender.ID)

	ufs.SentFriendRequest(c)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rec.Code)
	}
}

func TestGetFriendRequestsReturnsRequestsForReceiver(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db := testutil.SetupTestDB(t)

	req := models.FriendRequest{SenderID: "alice", ReceiverID: "bob"}
	if err := db.Create(&req).Error; err != nil {
		t.Fatalf("failed seeding request: %v", err)
	}

	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)
	c.Request = httptest.NewRequest(http.MethodGet, "/friends/requests", nil)
	c.Set("userID", "bob")

	ufs.GetFriendRequests(c)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rec.Code)
	}

	var response map[string]interface{}
	if err := json.Unmarshal(rec.Body.Bytes(), &response); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	data, ok := response["data"].([]interface{})
	if !ok || len(data) == 0 {
		t.Fatalf("expected non-empty data array in response, got %v", response["data"])
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
		"friendRequestId": reqModel.ID,
		"action":          "reject",
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

func TestFriendshipUniqueIndexPreventsDuplicates(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db := testutil.SetupTestDB(t)

	alice := models.User{Username: "alice_uniq", Password: "x"}
	bob := models.User{Username: "bob_uniq", Password: "x"}
	db.Create(&alice)
	db.Create(&bob)

	first := models.Friendship{UserID: alice.ID, FriendID: bob.ID, Status: "pending"}
	if err := db.Create(&first).Error; err != nil {
		t.Fatalf("first insert should succeed: %v", err)
	}

	duplicate := models.Friendship{UserID: alice.ID, FriendID: bob.ID, Status: "accepted"}
	if err := db.Create(&duplicate).Error; err == nil {
		t.Fatal("expected unique constraint violation on duplicate (user_id, friend_id), but got nil error")
	}
}

func TestSentFriendRequestRejectsSelf(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db := testutil.SetupTestDB(t)

	user := models.User{Username: "alice_self", Password: "x"}
	if err := db.Create(&user).Error; err != nil {
		t.Fatalf("failed seeding user: %v", err)
	}

	c, rec := newJSONContext(http.MethodPost, "/friends/request", gin.H{
		"friendId": user.ID,
	})
	c.Set("userID", user.ID)

	ufs.SentFriendRequest(c)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400 for self-request, got %d", rec.Code)
	}
}

func TestSentFriendRequestRejectsDuplicate(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db := testutil.SetupTestDB(t)

	alice := models.User{Username: "alice_dup", Password: "x"}
	bob := models.User{Username: "bob_dup", Password: "x"}
	db.Create(&alice)
	db.Create(&bob)

	db.Create(&models.Friendship{UserID: alice.ID, FriendID: bob.ID, Status: "pending"})

	c, rec := newJSONContext(http.MethodPost, "/friends/request", gin.H{
		"friendId": bob.ID,
	})
	c.Set("userID", alice.ID)

	ufs.SentFriendRequest(c)

	if rec.Code != http.StatusConflict {
		t.Fatalf("expected status 409 for duplicate request, got %d", rec.Code)
	}
}

func TestFriendRequestResponseRejectCleansUpFriendship(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db := testutil.SetupTestDB(t)

	alice := models.User{Username: "alice_rej", Password: "x"}
	bob := models.User{Username: "bob_rej", Password: "x"}
	db.Create(&alice)
	db.Create(&bob)

	reqModel := models.FriendRequest{SenderID: alice.ID, ReceiverID: bob.ID}
	db.Create(&reqModel)
	db.Create(&models.Friendship{UserID: alice.ID, FriendID: bob.ID, Status: "pending"})

	c, rec := newJSONContext(http.MethodPost, "/friends/respond", gin.H{
		"friendRequestId": reqModel.ID,
		"action":          "reject",
	})

	ufs.FriendRequestResponse(c)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rec.Code)
	}

	var remaining models.Friendship
	if err := db.Where("user_id = ? AND friend_id = ?", alice.ID, bob.ID).First(&remaining).Error; err == nil {
		t.Fatal("expected pending Friendship to be deleted after reject, but it still exists")
	}
}

func TestFriendRequestResponseAcceptCreatesBidirectionalFriendship(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db := testutil.SetupTestDB(t)

	alice := models.User{Username: "alice_bi", Password: "x"}
	bob := models.User{Username: "bob_bi", Password: "x"}
	db.Create(&alice)
	db.Create(&bob)

	reqModel := models.FriendRequest{SenderID: alice.ID, ReceiverID: bob.ID}
	db.Create(&reqModel)
	db.Create(&models.Friendship{UserID: alice.ID, FriendID: bob.ID, Status: "pending"})

	c, rec := newJSONContext(http.MethodPost, "/friends/respond", gin.H{
		"friendRequestId": reqModel.ID,
		"action":          "accept",
	})

	ufs.FriendRequestResponse(c)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rec.Code)
	}

	var fwd models.Friendship
	if err := db.Where("user_id = ? AND friend_id = ? AND status = ?", alice.ID, bob.ID, "accepted").First(&fwd).Error; err != nil {
		t.Fatalf("expected forward friendship (A→B accepted) to exist: %v", err)
	}

	var rev models.Friendship
	if err := db.Where("user_id = ? AND friend_id = ? AND status = ?", bob.ID, alice.ID, "accepted").First(&rev).Error; err != nil {
		t.Fatalf("expected reverse friendship (B→A accepted) to exist: %v", err)
	}
}
