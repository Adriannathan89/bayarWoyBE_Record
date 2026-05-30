package dto

type UpdateProfileDto struct {
	Username string `json:"username" binding:"required,min=3,max=50"`
}
