package dto

type RelationPhase2 struct {
	FriendRequestID string `json:"friendRequestId" binding:"required"`
	Action          string `json:"action" binding:"required,oneof=accept reject"`
}