package entity

import (
	"time"

	"github.com/google/uuid"
)

// OAuthCode represents an authorization code.
type OAuthCode struct {
	Id                  uuid.UUID     `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	Code                string        `gorm:"unique;not null" json:"code"`
	CodeChallenge       *string       `json:"codeChallenge"`
	CodeChallengeMethod string        `gorm:"not null" json:"codeChallengeMethod"`
	RedirectUri         *string       `json:"redirectUri"`
	UserId              *uuid.UUID    `gorm:"type:uuid" json:"userId"`
	User                *OAuthUser    `gorm:"foreignKey:UserID;references:ID" json:"user"`
	ClientId            *uuid.UUID    `gorm:"type:uuid" json:"clientId"`
	Client              *OAuthClient  `gorm:"foreignKey:ClientID;references:ID" json:"client"`
	ExpiresAt           time.Time     `gorm:"not null" json:"expiresAt"`
	CreatedAt           time.Time     `gorm:"autoCreateTime" json:"createdAt"`
	Scopes              []OAuthScope `gorm:"many2many:oauth_code_scopes;" json:"scopes"`
}

func (c *OAuthCode) GetScopeString() []string {
	var scopes []string
	for _, v := range c.Scopes {
		scopes = append(scopes, v.Name)
	}
	return scopes
}
