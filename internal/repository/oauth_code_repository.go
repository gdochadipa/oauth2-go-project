package repository

import (
	"context"

	"github.com/gdochadipa/oauth2-go-project/internal/entity"
)

type OAuthCodeRepostiory interface {
	FindByCode(ctx context.Context, code string) (*entity.OAuthCode, error)
	IssueWithAuthCode(ctx context.Context, client *entity.OAuthClient, user *entity.OAuthUser, scope []string) (*entity.OAuthCode, error)
	PersistOAuthCode(ctx context.Context, authCode *entity.OAuthCode) error
	IsRevokedOAuthCode(ctx context.Context, code string) (bool, error)
	RevokedOAuthCode(ctx context.Context, code string) error
	GetByRefreshToken(ctx context.Context, refreshTokenCode string) (*entity.OAuthToken, error)
}

// FindByCode implements Repository.
func (r *dbRepository) FindByCode(ctx context.Context, code string) (*entity.OAuthCode, error) {
	panic("unimplemented")
}

// GetByRefreshToken implements Repository.
func (r *dbRepository) GetByRefreshToken(ctx context.Context, refreshTokenCode string) (*entity.OAuthToken, error) {
	panic("unimplemented")
}

// IsRevokedOAuthCode implements Repository.
func (r *dbRepository) IsRevokedOAuthCode(ctx context.Context, code string) (bool, error) {
	panic("unimplemented")
}

// IssueWithAuthCode implements Repository.
func (r *dbRepository) IssueWithAuthCode(ctx context.Context, client *entity.OAuthClient, user *entity.OAuthUser, scope []string) (*entity.OAuthCode, error) {
	panic("unimplemented")
}

// PersistOAuthCode implements Repository.
func (r *dbRepository) PersistOAuthCode(ctx context.Context, authCode *entity.OAuthCode) error {
	panic("unimplemented")
}

// RevokedOAuthCode implements Repository.
func (r *dbRepository) RevokedOAuthCode(ctx context.Context, code string) error {
	panic("unimplemented")
}
