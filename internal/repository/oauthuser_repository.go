package repository

import (
	"fmt"

	"github.com/gdochadipa/oauth2-go-project/internal/entity"
	"github.com/google/uuid"
)

func (db *gormDB) GetUser(id string) (*entity.OAuthUser, error) {
	var user entity.OAuthUser
	parsedUUID, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("Something wrong with ID")
	}

	result := db.db.First(&user, "id = ?", parsedUUID)
	if result.Error != nil {
		return nil, result.Error
	}

	return &user, nil
}
