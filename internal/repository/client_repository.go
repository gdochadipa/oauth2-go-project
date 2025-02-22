package repository

import (
	"context"

	"github.com/gdochadipa/oauth2-go-project/internal/entity"
	"github.com/gdochadipa/oauth2-go-project/internal/enum"
	"github.com/google/uuid"
)

type ClientRepository interface {
	OAuthClientFindById(ctx context.Context, clientId uuid.UUID) (*entity.OAuthClient, error)
	IsClientValid(ctx context.Context, grantType *enum.GrantIdentifier, client *entity.OAuthClient, clientSecret *string) bool
}

// FindById implements Repository.
func (r *dbRepository) OAuthClientFindById(ctx context.Context, clientId uuid.UUID) (*entity.OAuthClient, error) {
	panic("unimplemented")
}

// IsClientValid implements Repository.
func (r *dbRepository) IsClientValid(ctx context.Context, grantType *enum.GrantIdentifier, client *entity.OAuthClient, clientSecret *string) bool {
	panic("unimplemented")
}
