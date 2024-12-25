package entity

import "github.com/google/uuid"

type Item struct {
	id          uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()`
	name        string    `gorm:"type:varchar(100);not null"`
	description string    `gorm:"type:varchar(20);unique;not null"`
}
