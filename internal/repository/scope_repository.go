package repository

import (
	"context"

	"github.com/gdochadipa/oauth2-go-project/internal/entity"
	"github.com/gdochadipa/oauth2-go-project/internal/enum"
	"github.com/google/uuid"
)

type ScopeRepository interface {
	GetAllScope(ctx context.Context, scopeNames []string) ([]entity.OAuthScope, error)
	GetScopeByIdentifier(ctx context.Context, grant enum.GrantIdentifier, scopes []entity.OAuthScope, client *entity.OAuthClient, userID *uuid.UUID) ([]entity.OAuthScope, error)
}

// GetAllScope implements Repository.
// like this verify is scope already is exists on our system
func (r *dbRepository) GetAllScope(ctx context.Context, scopeNames []string) ([]entity.OAuthScope, error) {
	panic("unimplemented")
}

// get scope is exists in our system
// by grant, clientId and userId ?
// this function also comparing scope on our system with scope by token/request
func (r *dbRepository) GetScopeByIdentifier(ctx context.Context, grant enum.GrantIdentifier, scopes []entity.OAuthScope, client *entity.OAuthClient, userID *uuid.UUID) ([]entity.OAuthScope, error) {
	panic("unimplemented")
}
