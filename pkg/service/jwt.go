package service

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type JWTInterface interface {
	createRefreshToken(userId uuid.UUID, clientId string, scope []string, expiredAt time.Time, issueAt time.Time) (*string, error)
	createAccessToken(userId uuid.UUID, username string, scope []string, expiredAt time.Time, issueAt time.Time, notBefore time.Time) (*string, error)
	decodeAccessToken(encryptData string) (*AccessToken, error)
	decodeRefreshToken(encryptData string) (*RefreshAccessToken, error)
	verifyAccessToken(encryptData string) (*AccessToken, error)
	verifyRefreshToken(encryptData string) (*RefreshAccessToken, error)
}

type JWTRepository struct {
	secretKey []byte
}

/*
*
ada dua data
AccessToken

# RefreshToken

keduanya beda isinya dan fungsionalnya. bisa di atur sih seharusnya, cuma cari lebih dalam soal isi RefreshToken,
sepertinya bisa sama dengan AccessToken tapi isinya lebih banyak aja.

create AccessToken akan beda dg RefreshToken, yg memedakan cuma issian aja, dan parameter-
*/
type AccessToken struct {
	UserId      string   `json:"userId"`
	Username    string   `json:"username"`
	ClientId    string   `json:"clientId"`
	RedirectUri string   `json:"redirectUri"`
	AuthCodeId  string   `json:"auth_code"`
	Scopes      []string `json:"scopes"`
	jwt.RegisteredClaims
}

type RefreshAccessToken struct {
	ClientId            string    `json:"clientId"`
	RedirectUri         string    `json:"redirectUri"`
	AuthCodeId          string    `json:"auth_code"`
	Scopes              []string  `json:"scopes"`
	UserId              string    `json:"userId"`
	ExpireTime          time.Time `json:"expireTime"`
	CodeChallenge       *string   `json:"codeChallenge"`
	CodeChallengeMethod *string   `json:"codeChallengeMethod"`
	jwt.RegisteredClaims
}

func (r *JWTRepository) createRefreshToken(userId uuid.UUID, clientId string, scope []string, expiredAt time.Time, issueAt time.Time) (*string, error) {
	claims := RefreshAccessToken{
		UserId:   userId.String(),
		ClientId: clientId,
		Scopes:   scope,
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

func (r *JWTRepository) createAccessToken(userId uuid.UUID, username string, scope []string, expiredAt time.Time, issueAt time.Time, notBefore time.Time) (*string, error) {
	claims := AccessToken{
		UserId:   userId.String(),
		Username: username,
		Scopes:   scope,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiredAt),
			Audience:  []string{userId.String(), username},
			Subject:   userId.String(),
			IssuedAt:  jwt.NewNumericDate(issueAt),
			NotBefore: jwt.NewNumericDate(notBefore),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString(r.secretKey)
	if err != nil {
		return nil, fmt.Errorf("error signing token: %w", err)
	}

	return &signedToken, nil
}

func (r *JWTRepository) decodeAccessToken(encryptData string) (*AccessToken, error) {
	token, _ := jwt.ParseWithClaims(encryptData, &AccessToken{}, func(t *jwt.Token) (interface{}, error) {
		return nil, nil
	})

	if claims, ok := token.Claims.(*AccessToken); ok {
		return claims, nil
	}

	return nil, fmt.Errorf("invalid token claims")
}

func (r *JWTRepository) decodeRefreshToken(encryptData string) (*RefreshAccessToken, error) {
	token, _ := jwt.ParseWithClaims(encryptData, &RefreshAccessToken{}, func(t *jwt.Token) (interface{}, error) {
		return nil, nil
	})

	if claims, ok := token.Claims.(*RefreshAccessToken); ok {
		return claims, nil
	}

	return nil, fmt.Errorf("invalid refresh token claims")
}

func (r *JWTRepository) verifyAccessToken(encryptData string) (*AccessToken, error) {
	token, err := jwt.ParseWithClaims(encryptData, &AccessToken{}, func(t *jwt.Token) (interface{}, error) {

		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return r.secretKey, nil
	})

	if err != nil {
		return nil, fmt.Errorf("error parsing token: %w", err)
	}

	if claims, ok := token.Claims.(*AccessToken); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("invalid token claims")
}

func (r *JWTRepository) verifyRefreshToken(encryptData string) (*RefreshAccessToken, error) {
	token, err := jwt.ParseWithClaims(encryptData, &RefreshAccessToken{}, func(t *jwt.Token) (interface{}, error) {

		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return r.secretKey, nil
	})

	if err != nil {
		return nil, fmt.Errorf("error parsing refresh token: %w", err)
	}

	if claims, ok := token.Claims.(*RefreshAccessToken); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("invalid refresh token claims")
}

func NewJWTRepository(secretKey []byte) JWTInterface {
	return &JWTRepository{secretKey: secretKey}
}
