package repository

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type JWTRepository struct {
	secretKey []byte
}

type AccessToken struct {
	UserId   string `json:"userId"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

type RefreshAccessToken struct {
	jwt.RegisteredClaims
}

func (r *JWTRepository) createAccessToken(userId uuid.UUID, username string) (*string, error) {
	claims := AccessToken{
		UserId:   userId.String(),
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString(r.secretKey)
	if err != nil {
		return nil, fmt.Errorf("error signing token: %w", err)
	}

	return &signedToken, nil
}
