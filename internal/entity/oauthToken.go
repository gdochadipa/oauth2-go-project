package entity

import (
	"time"

	"github.com/google/uuid"
)

// OAuthToken represents an access or refresh token.
type OAuthToken struct {
	Id                    uuid.UUID     `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	AccessToken           string        `gorm:"unique;not null" json:"accessToken"`
	AccessTokenExpiresAt  time.Time     `gorm:"not null" json:"accessTokenExpiresAt"`
	RefreshToken          *string       `gorm:"unique" json:"refreshToken"`
	RefreshTokenExpiresAt *time.Time    `json:"refreshTokenExpiresAt"`
	ClientId              *uuid.UUID    `gorm:"type:uuid" json:"clientId"`
	Client                *OAuthClient  `gorm:"foreignKey:ClientID;references:ID" json:"client"`
	UserId                *uuid.UUID    `gorm:"type:uuid" json:"userId"`
	User                  *OAuthUser    `gorm:"foreignKey:UserID;references:ID" json:"user"`
	Scopes                []*OAuthScope `gorm:"many2many:oauth_token_scopes;" json:"scopes"`
	CreatedAt             time.Time     `gorm:"autoCreateTime" json:"createdAt"`
}
