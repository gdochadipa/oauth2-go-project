package repository

import (
	"context"

	"github.com/gdochadipa/oauth2-go-project/internal/entity"
)

type ScopeRepository interface {
	GetAllScope(ctx context.Context, scopeNames []string) ([]entity.OAuthScope, error)
}

// GetAllScope implements Repository.
func (r *dbRepository) GetAllScope(ctx context.Context, scopeNames []string) ([]entity.OAuthScope, error) {
	panic("unimplemented")
}
