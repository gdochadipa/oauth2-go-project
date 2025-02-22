package service

import (
	"context"
	"fmt"
	"time"

	"github.com/gdochadipa/oauth2-go-project/internal/entity"
	"github.com/gdochadipa/oauth2-go-project/internal/enum"
	"github.com/google/uuid"
	"google.golang.org/grpc/metadata"
)

type AuthorizationRequest struct {
	Scopes              []entity.OAuthScope
	IsAuthApproved      bool
	RedirectUri         *string
	State               *string
	CodeChallenge       *string
	CodeChallengeMethod *enum.CodeEnum
	Client              *entity.OAuthClient
	GrantTypeId         *enum.GrantIdentifier
	User                *entity.OAuthUser
	Audience            []string
}

type PayloadAuthenticationCode struct {
	ClidenID            string         `json:"client_id"`
	AuthCodeID          string         `json:"auth_code_id"`
	ExpireTime          int64          `json:"expire_time"`
	Scopes              []string       `json:"scopes"`
	UserID              *int16         `json:"user_id"`
	RedirectURI         *string        `json:"redirect_uri"`
	CodeChallenge       *string        `json:"code_challenge"`
	CodeChallengeMethod *enum.CodeEnum `json:"code_challenge_method"`
	Audience            []string       `json:"audience"`
}

type AuthCodeGrantInterface interface {
	RespondToAccessTokenRequest(ctx context.Context, expiredAccessToken *time.Time, metadata metadata.MD) (*BearerTokenResponse, error)
	ValidateAuthorizationRequest(ctx context.Context, clientId *string, scopes []entity.OAuthScope, state *string, audience *string, codeChallenge *string, codeChallengeMethod *string) (*AuthorizationRequest, error)
	CompleteAuthorizationRequest(ctx context.Context, authRequest *AuthorizationRequest) (*string, error)
	IssueAuthCode(ctx context.Context, expiredAuthCode *time.Time, client *entity.OAuthClient, userId *uuid.UUID, redirectUri *string, codeChallenge *string, codeChallengeMethod *string, scopes []entity.OAuthScope) (*entity.OAuthCode, error)
	RespondToRevokeRequest(ctx context.Context, grantType *string, token *string) error
}

// CompleteAuthorizationRequest implements ServiceInterface.
func (g *ServiceServer) CompleteAuthorizationRequest(ctx context.Context, authRequest *AuthorizationRequest) (*string, error) {

	if authRequest == nil {
		return nil, fmt.Errorf("Invalid authorization request")
	}

	if !authRequest.IsAuthApproved {
		return nil, fmt.Errorf("Auth is not approved")
	}

	if authRequest.User == nil {
		return nil, fmt.Errorf("User is not found in request")
	}

	expired := g.dateInterval.GetEndDate()
	authCode, error := g.IssueAuthCode(ctx, &expired, authRequest.Client, &authRequest.User.Id, authRequest.RedirectUri, authRequest.CodeChallenge, (*string)(authRequest.CodeChallengeMethod), authRequest.Scopes)

	if error != nil {
		return nil, error
	}

	payload := PayloadAuthenticationCode{
		ClidenID:            authCode.Client.Id.String(),
		RedirectURI:         authCode.RedirectUri,
		AuthCodeID:          authCode.Code,
		Scopes:              authCode.Scopes,
		UserID:              &authCode.User.IdentifierId,
		ExpireTime:          g.dateInterval.GetEndTimeSeconds(),
		CodeChallenge:       authRequest.CodeChallenge,
		CodeChallengeMethod: authRequest.CodeChallengeMethod,
		Audience:            authRequest.Audience,
	}

	authCodeToken, err := g.jwt.createAuthCodeToken(&payload)

	if err != nil {
		return nil, err
	}

	redirectUri, err := g.makeRedirectUrl(authRequest.RedirectUri, map[string]interface{}{"code": authCodeToken})

	if err != nil {
		return nil, err
	}
	//! belum selesai
	// perlu ngedefine redirect response sepertinya

	return redirectUri, nil
}

// IssueAuthCode implements ServiceInterface.
func (g *ServiceServer) IssueAuthCode(ctx context.Context, expiredAuthCode *time.Time, client *entity.OAuthClient, userId *uuid.UUID, redirectUri *string, codeChallenge *string, codeChallengeMethod *string, scopes []entity.OAuthScope) (*entity.OAuthCode, error) {
	panic("unimplemented")
}

// RespondToAccessTokenRequest implements ServiceInterface.
func (g *ServiceServer) RespondToAccessTokenRequest(ctx context.Context, expiredAccessToken *time.Time, metadata metadata.MD) (*BearerTokenResponse, error) {
	panic("unimplemented")
}

// RespondToRevokeRequest implements ServiceInterface.
func (g *ServiceServer) RespondToRevokeRequest(ctx context.Context, grantType *string, token *string) error {
	panic("unimplemented")
}

// ValidateAuthorizationRequest implements ServiceInterface.
func (g *ServiceServer) ValidateAuthorizationRequest(ctx context.Context, clientId *string, scopes []entity.OAuthScope, state *string, audience *string, codeChallenge *string, codeChallengeMethod *string) (*AuthorizationRequest, error) {
	panic("unimplemented")
}
