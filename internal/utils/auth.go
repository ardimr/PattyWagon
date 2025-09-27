package utils

import (
	"WeenieHut/internal/constants"
	"WeenieHut/internal/model"
	"context"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

var jwtSignatureKey string

func init() {
	jwtSignatureKey = os.Getenv("JWT_SIGNATURE_KEY")
}

// GetUserIDFromCtx get user ID from ctx
func GetUserIDFromCtx(ctx context.Context) (int64, bool) {
	userID, ok := ctx.Value(constants.UserIDCtxKey).(int64)
	return userID, ok
}

// HashPassword generates a bcrypt hash for the given password.
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 4)
	return string(bytes), err
}

// VerifyPassword verifies if the given password matches the stored hash.
func VerifyPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// GenerateToken generate jwt token with userID information
func GenerateToken(userID int64) (string, error) {
	claims := model.Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24)),
			Issuer:    "WeenieHut",
		},
		UserID: userID,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, err := token.SignedString([]byte(jwtSignatureKey))
	if err != nil {
		return "", err
	}
	return ss, nil
}

// ParseUserIDFromToken verify jwt token and get user ID information
func ParseUserIDFromToken(tokenString string) (int64, error) {
	token, err := jwt.ParseWithClaims(tokenString, &model.Claims{}, func(token *jwt.Token) (any, error) {
		return []byte(jwtSignatureKey), nil
	})
	if err != nil {
		return 0, err
	}
	claims, ok := token.Claims.(*model.Claims)
	if !ok || claims.UserID == 0 {
		return 0, err
	}
	return claims.UserID, nil
}
