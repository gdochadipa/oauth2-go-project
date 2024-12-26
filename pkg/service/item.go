package service

import (
	"context"

	"github.com/gdochadipa/oauth2-go-project/internal/entity"
	"github.com/google/uuid"
)

type ItemService interface {
	PostItem(ctx context.Context, name string, descriptiont string) (*entity.Item, error)
}

func (s *oauthService) PostItem(ctx context.Context, name string, desc string) (*entity.Item, error) {
	item := &entity.Item{
		Name:        name,
		Description: desc,
	}

	newItem, err := s.repository.InsertItem(ctx, *item)
	if err != nil {
		return nil, err
	}

	return newItem, nil
}

func (s *oauthService) GetAllItems(ctx context.Context, skip uint32, limitPage uint32) ([]entity.Item, error) {
	if limitPage > 15 || (skip == 0 && limitPage == 0) {
		limitPage = 15
	}

	return s.repository.GetListItem(ctx, skip, limitPage)
}

func (s *oauthService) GetItem(ctx context.Context, id uuid.UUID) (*entity.Item, error) {
	return s.repository.GetItem(ctx, id)
}

func (s *oauthService) PutItem(ctx context.Context, item entity.Item) (*entity.Item, error) {
	return s.repository.UpdateItem(ctx, item.Id, item)
}

func (s *oauthService) DropItem(ctx context.Context, id uuid.UUID) error {
	if err := s.repository.DeleteItem(ctx, id); err != nil {
		return err
	}

	return nil
}
