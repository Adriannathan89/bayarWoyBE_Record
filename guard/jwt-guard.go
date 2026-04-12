package guard

import (
	"bayar-woy-project/config"
	"time"

	"bayar-woy-project/models"

	"github.com/golang-jwt/jwt/v5"
)

var jwtKey = []byte(config.GetEnv("JWT_SECRET_KEY"))
var refreshJwtKey = []byte(config.GetEnv("JWT_REFRESH_SECRET_KEY"))

func GenerateToken(username string, userID string) (string, error) {
	expirationTime := time.Now().Add(15 * time.Minute)

	claims := models.Claims{
		Username: username,
		UserID:   userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			Issuer:    "my-app",
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString(jwtKey)
}

func GenerateRefreshToken(username string, userID string) (string, error) {
	expirationTime := time.Now().Add(4 * time.Hour)
	claims := models.Claims{
		Username: username,
		UserID:   userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			Issuer:    "my-app",
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return refreshToken.SignedString(refreshJwtKey)
}

func ValidateToken(tokenString string) (*models.Claims, error) {
	claims := &models.Claims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	if(err != nil || !token.Valid) {	
		return nil, err
	}

	return claims, nil
}