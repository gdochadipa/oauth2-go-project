package service

import (
	"context"
	"time"

	"github.com/gdochadipa/oauth2-go-project/internal/entity"
)

type TokenExchangeGrantInterface interface {
	RespondToExchangeRequest(ctx context.Context, expiredAccessToken *time.Time, subjectToken *string, subjectTokenType *string, actorToken *string, actorTokenType *string, scopes []entity.OAuthScope) (*BearerTokenResponse, error)
}

// RespondToExchangeRequest implements ServiceInterface.
func (g *ServiceServer) RespondToExchangeRequest(ctx context.Context, expiredAccessToken *time.Time, subjectToken *string, subjectTokenType *string, actorToken *string, actorTokenType *string, scopes []entity.OAuthScope) (*BearerTokenResponse, error) {
	panic("unimplemented")
}
