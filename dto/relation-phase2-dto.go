package dto

type RelationPhase2 struct {
	FriendshipID string `json:"friendshipId" binding:"required"`
	Action       string `json:"action" binding:"required,oneof=accept reject"`
}