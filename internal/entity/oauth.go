package entity

import (
	"time"

	"github.com/gdochadipa/oauth2-go-project/internal/enum"
	"github.com/google/uuid"
)

type OAuthUser struct {
	Id           uuid.UUID `json:"id"`
	Username     string    `json:"username"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"password"`
	CreatedAt    time.Time `json:"CreatedAt"`
	UpdatedAt    time.Time `json:"UpdatedAt"`
}

type OAuthCode struct {
	Id                  uuid.UUID     `json:"id"`
	Code                string        `json:"code"`
	CodeChallenge       *string       `json:"codeChallenge"`
	CodeChallengeMethod enum.CodeEnum `json:"codeChallengeMethod"`
	RedirectUri         *string       `json:"redirectUri"`
	UserId              *uuid.UUID    `json:"userId"`
	User                *OAuthUser    `json:"User"`
	ClientId            *uuid.UUID    `json:"clientId"`
	Client              *OAuthClient  `json:"Client"`
	ExpairedAt          time.Time     `json:"expairedAt"`
	CreatedAt           time.Time     `json:"createdAt"`
	Scopes              []string      `json:"Scopes"`
}

type OAuthClient struct {
	Id            uuid.UUID            `json:"id"`
	Name          string               `json:"name"`
	RedirectUris  []string             `json:"redirectUris"`
	Secret        *string              `json:"secret"`
	ClientId      *uuid.UUID           `json:"clientId"`
	AllowedGrants enum.GrantIdentifier `json:"allowedGrants"`
	CreatedAt     time.Time            `json:"createdAt"`
	Scopes        []string             `json:"Scopes"`
}

type OAuthToken struct {
	AccessToken           string     `json:"accessToken"`
	AccessTokenExpiresAt  time.Time  `json:"accessTokenExpiresAt"`
	RefreshToken          *string    `json:"refreshToken"`
	RefreshTokenExpiresAt *time.Time `json:"refreshTokenExpiresAt"`
	ClientId              *uuid.UUID `json:"clientId"`
	UserId                *uuid.UUID `json:"userId"`
	Scopes                []string   `json:"scopes"`
	CreatedAt             time.Time  `json:"createdAt"`
}
