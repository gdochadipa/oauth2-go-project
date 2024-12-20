package repository

import "gorm.io/gorm"

type gormDB struct {
	db *gorm.DB
}
