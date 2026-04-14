package dto

type RelationPhase1 struct {
	FriendID string `json:"friendId" binding:"required"`
}