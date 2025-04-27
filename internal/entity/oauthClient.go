package entity

import (
	"time"

	"github.com/gdochadipa/oauth2-go-project/internal/enum"
	"github.com/google/uuid"
)

// OAuthClient represents an OAuth client.
type OAuthClient struct {
	Id            uuid.UUID            `gorm:"column:Id;type:uuid;primaryKey;" json:"id"`
	Name          string               `gorm:"column:Name;not null;" json:"name"`
	RedirectUris  []string             `gorm:"column:RedirectUris;type:text[]" json:"redirectUris"`
	Secret        *string              `gorm:"column:Secret;type:varchar(255)" json:"secret"`
	ClientId      *uuid.UUID           `gorm:"column:ClientId;type:uuid;unique" json:"clientId"`
	AllowedGrants enum.GrantIdentifier `gorm:"column:AllowedGrants;not null" json:"allowedGrants"`
	CreatedAt     time.Time            `gorm:"column:CreatedAt;autoCreateTime" json:"createdAt"`
	Scopes        []OAuthScope         `gorm:"column:Scopes;many2many:oauth_client_scopes;" json:"scopes"`
}
