package entity

import "github.com/google/uuid"

type OAuthUser struct {
	Id       uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()`
	Username string    `gorm:"type:varchar(100);not null"`
	email    string    `gorm:"type:varchar(20);unique;not null"`
}
