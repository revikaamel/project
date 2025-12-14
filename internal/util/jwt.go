package util

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"uas-backend/internal/model"
)

func GenerateToken(userID, role, email, secret string, expireHours int) (string, error) {
	claims := &model.TokenClaims{
		UserID: userID,
		Email:  email,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(expireHours) * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(secret))
}
