package dto

type ValidateOtpDto struct {
	UserID string `json:"userId" binding:"required"`
	OTP    string `json:"otp" binding:"required"`
}