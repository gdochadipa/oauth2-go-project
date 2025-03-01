package repository

import (
	"context"

	"github.com/gdochadipa/oauth2-go-project/internal/entity"
	"github.com/google/uuid"
)

type UserRepository interface {
	GetUserByIdentifier(ctx context.Context, identifier *uuid.UUID) (*entity.OAuthUser, error)
}

// GetUserByIdentifiry implements Repository.
func (r *dbRepository) GetUserByIdentifier(ctx context.Context, identifier *uuid.UUID) (*entity.OAuthUser, error) {
	panic("unimplemented")
}
