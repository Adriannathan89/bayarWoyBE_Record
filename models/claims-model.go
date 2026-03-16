package models

import "github.com/golang-jwt/jwt/v5"

type Claims struct {
	UserID   string `json:"userId"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}
