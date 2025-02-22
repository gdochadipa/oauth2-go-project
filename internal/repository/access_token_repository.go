package repository

import (
	"context"

	"github.com/gdochadipa/oauth2-go-project/internal/entity"
)

type AccessTokenRepository interface {
	IssueToken(ctx context.Context, client *entity.OAuthClient, scope []entity.OAuthScope, user *entity.OAuthUser) (*entity.OAuthToken, error)
	IssueRefreshToken(ctx context.Context, accessToken *entity.OAuthToken, client *entity.OAuthClient) (*entity.OAuthToken, error)
	PersistAccessToken(ctx context.Context, accessToken *entity.OAuthToken) error
	RevokeAccessToken(ctx context.Context, accessToken *entity.OAuthToken) error
	IsRefreshTokenRevoked(ctx context.Context, refreshToken *entity.OAuthToken) (bool, error)
}

// IsRefreshTokenRevoked implements Repository.
func (r *dbRepository) IsRefreshTokenRevoked(ctx context.Context, refreshToken *entity.OAuthToken) (bool, error) {
	panic("unimplemented")
}

// IssueRefreshToken implements Repository.
func (r *dbRepository) IssueRefreshToken(ctx context.Context, accessToken *entity.OAuthToken, client *entity.OAuthClient) (*entity.OAuthToken, error) {
	panic("unimplemented")
}

// IssueToken implements Repository.
func (r *dbRepository) IssueToken(ctx context.Context, client *entity.OAuthClient, scope []entity.OAuthScope, user *entity.OAuthUser) (*entity.OAuthToken, error) {
	panic("unimplemented")
}

// PersistAccessToken implements Repository.
func (r *dbRepository) PersistAccessToken(ctx context.Context, accessToken *entity.OAuthToken) error {
	panic("unimplemented")
}

// RevokeAccessTOken implements Repository.
func (r *dbRepository) RevokeAccessToken(ctx context.Context, accessToken *entity.OAuthToken) error {
	panic("unimplemented")
}
