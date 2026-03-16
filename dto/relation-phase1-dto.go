package dto

type RelationPhase1 struct {
	FriendUsername string `json:"friendUsername" binding:"required"`
}