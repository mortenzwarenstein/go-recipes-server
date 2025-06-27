package util

import (
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

type CustomClaims struct {
	jwt.RegisteredClaims
	Email string `json:"email"`
}

var (
	AccessTokenSecret  = []byte("secret")  // TODO: Replace with env
	RefreshTokenSecret = []byte("refresh") // TODO: Replace with env
	AccessExpiryTime   = time.Minute * 10
	RefreshExpiryTime  = time.Hour * 24 * 7
)

func GenerateTokens(userID string, email string) (accessToken string, refreshToken string, err error) {
	accessClaims := CustomClaims{
		jwt.RegisteredClaims{
			Subject:   userID,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(AccessExpiryTime)), // TODO: Replace with env
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "go-recipe-api", // TODO: Replace with env
		},
		email,
	}

	at := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	accessToken, err = at.SignedString(AccessTokenSecret)
	if err != nil {
		return
	}

	refreshClaims := accessClaims
	refreshClaims.ExpiresAt = jwt.NewNumericDate(time.Now().Add(RefreshExpiryTime)) // TODO: Replace with env
	refreshToken, err = at.SignedString(RefreshTokenSecret)
	return
}

func VerifyAccessToken(tokenStr string) (*CustomClaims, error) {
	return verifyToken(tokenStr, AccessTokenSecret)
}

func VerifyRefreshToken(tokenStr string) (*CustomClaims, error) {
	return verifyToken(tokenStr, RefreshTokenSecret)
}

func verifyToken(tokenStr string, secret []byte) (*CustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return secret, nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*CustomClaims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}
