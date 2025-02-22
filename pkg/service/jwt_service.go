package service

import (
	"fmt"
	"time"

	"github.com/gdochadipa/oauth2-go-project/internal/enum"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type JWTInterface interface {
	createRefreshToken(userId uuid.UUID, clientId string, scopes string, expiredAt time.Time) (*string, error)
	createAccessToken(userId uuid.UUID, username string, scope string, clientId string, expiredAt time.Time, redirectUri string, issueAt time.Time, notBefore time.Time) (*string, error)
	decodeAccessToken(encryptData string) (*AccessToken, error)
	decodeRefreshToken(encryptData string) (*RefreshAccessToken, error)
	verifyAccessToken(encryptData string) (*AccessToken, error)
	verifyRefreshToken(encryptData string) (*RefreshAccessToken, error)
	createAuthCodeToken(payload *PayloadAuthenticationCode) (*string, error)
	verifyAuthCodeToken(tokenString string) (*AuthCodeToken, error)
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
	UserId      string `json:"userId"`
	Username    string `json:"username"`
	ClientId    string `json:"clientId"`
	RedirectUri string `json:"redirectUri"`
	AuthCodeId  string `json:"auth_code"`
	Scopes      string `json:"scopes"`
	jwt.RegisteredClaims
}

type RefreshAccessToken struct {
	ClientId            string    `json:"clientId"`
	RedirectUri         string    `json:"redirectUri"`
	AuthCodeId          string    `json:"auth_code"`
	Scopes              string    `json:"scopes"`
	UserId              string    `json:"userId"`
	ExpireTime          time.Time `json:"expireTime"`
	CodeChallenge       *string   `json:"codeChallenge"`
	CodeChallengeMethod *string   `json:"codeChallengeMethod"`
	jwt.RegisteredClaims
}

type AuthCodeToken struct {
	ClidenID            string         `json:"client_id"`
	AuthCodeID          string         `json:"auth_code_id"`
	ExpireTime          int64          `json:"expire_time"`
	Scopes              []string       `json:"scopes"`
	UserID              *int16         `json:"user_id"`
	RedirectURI         *string        `json:"redirect_uri"`
	CodeChallenge       *string        `json:"code_challenge"`
	CodeChallengeMethod *enum.CodeEnum `json:"code_challenge_method"`
	Audience            []string       `json:"audience"`
	jwt.RegisteredClaims
}

func (r *JWTRepository) createRefreshToken(userId uuid.UUID, clientId string, scopes string, expiredAt time.Time) (*string, error) {
	claims := RefreshAccessToken{
		UserId:   userId.String(),
		ClientId: clientId,
		Scopes:   scopes,
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

func (r *JWTRepository) createAccessToken(userId uuid.UUID, username string, scope string, clientId string, expiredAt time.Time, redirectUri string, issueAt time.Time, notBefore time.Time) (*string, error) {
	claims := AccessToken{
		UserId:      userId.String(),
		Username:    username,
		Scopes:      scope,
		ClientId:    clientId,
		RedirectUri: redirectUri,
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

func (r *JWTRepository) createAuthCodeToken(payload *PayloadAuthenticationCode) (*string, error) {
	claims := AuthCodeToken{
		ClidenID:            payload.ClidenID,
		RedirectURI:         payload.RedirectURI,
		AuthCodeID:          payload.AuthCodeID,
		Scopes:              payload.Scopes,
		UserID:              payload.UserID,
		ExpireTime:          payload.ExpireTime,
		CodeChallenge:       payload.CodeChallenge,
		CodeChallengeMethod: payload.CodeChallengeMethod,
		Audience:            payload.Audience,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)), // Token expires in 24 hours
			IssuedAt:  jwt.NewNumericDate(time.Now()),                     // Token issue time
			NotBefore: jwt.NewNumericDate(time.Now()),                     // Token not valid before this time
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	code, err := token.SignedString(r.secretKey)

	if err != nil {
		return nil, fmt.Errorf("error parsing refresh token: %w", err)
	}

	return &code, nil
}

func (r *JWTRepository) verifyAuthCodeToken(tokenString string) (*AuthCodeToken, error) {
	token, err := jwt.ParseWithClaims(tokenString, &AuthCodeToken{}, func(t *jwt.Token) (interface{}, error) {

		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return r.secretKey, nil
	})

	if err != nil {
		return nil, fmt.Errorf("error parsing token: %v", err)
	}

	// Extract and validate claims
	if claims, ok := token.Claims.(*AuthCodeToken); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("invalid token claims")
}

func NewJWTRepository(secretKey []byte) JWTInterface {
	return &JWTRepository{secretKey: secretKey}
}
