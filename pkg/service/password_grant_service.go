package service

import (
	"context"

	"github.com/gdochadipa/oauth2-go-project/internal/entity"
)

type PasswordGrantInterface interface {
	respondToAccessTokenRequest(ctx context.Context, username string, password string, grantType string, scopes []entity.OAuthScope) (*BearerTokenResponse, error)
	validateUser(ctx context.Context, username string, password string, client *entity.OAuthClient) (*entity.OAuthUser, error)
}

// respondToAccessTokenRequest implements ServiceInterface.
// metadata metadata.MD also check metadata
// you can get metadata from context metadata.FromIncomingContext(ctx)
func (g *ServiceServer) respondToAccessTokenRequest(ctx context.Context, username string, password string, grantType string, scopes []entity.OAuthScope) (*BearerTokenResponse, error) {
	panic("unimplemented")
}

// validateUser implements ServiceInterface.
func (g *ServiceServer) validateUser(ctx context.Context, username string, password string, client *entity.OAuthClient) (*entity.OAuthUser, error) {
	panic("unimplemented")
}
