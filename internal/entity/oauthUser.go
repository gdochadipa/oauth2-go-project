package entity

import (
	"time"

	"github.com/google/uuid"
)

// OAuthUser represents a user in the OAuth system.
type OAuthUser struct {
	Id           uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	IdentifierId int16     `gorm:"not null" json:"identifierId"`
	Username     string    `gorm:"unique;not null" json:"username"`
	Email        string    `gorm:"unique;not null" json:"email"`
	PasswordHash string    `gorm:"not null" json:"password"`
	CreatedAt    time.Time `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt    time.Time `gorm:"autoUpdateTime" json:"updatedAt"`
}
