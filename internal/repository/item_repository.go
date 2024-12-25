package repository

import (
	"fmt"

	"github.com/gdochadipa/oauth2-go-project/internal/entity"
	"github.com/google/uuid"
)

func (r *Repository) GetItemDetail(id string) (*entity.Item, error) {
	var item entity.Item
	parseUUID, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("something wrong with id")
	}

	result := r.db.First(&item, "id = ?", parseUUID)
	if result.Error != nil {
		return nil, result.Error
	}
	return &item, err
}
