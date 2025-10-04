package model

import (
	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	jwt.RegisteredClaims
	UserID int64 `json:"userID"`
	Role   int16 `json:"role"`
}
