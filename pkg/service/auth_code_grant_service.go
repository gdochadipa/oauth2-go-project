package service

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/gdochadipa/oauth2-go-project/internal/entity"
	"github.com/gdochadipa/oauth2-go-project/internal/enum"
	"github.com/google/uuid"
	"google.golang.org/grpc/metadata"
)

type AuthorizationRequest struct {
	Scopes              *string
	IsAuthApproved      bool
	RedirectUri         *string
	State               *string
	CodeChallenge       *string
	CodeChallengeMethod *enum.CodeEnum
	Client              *entity.OAuthClient
	GrantTypeId         *enum.GrantIdentifier
	User                *entity.OAuthUser
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

type AuthCodeRequest struct {
	grantType     string
	authorization string
	redirectUri   *string
	code          *string
	clientId      *string
	codeVerified  *string
}

type AuthCodeGrantInterface interface {
	RespondToAccessTokenRequest(ctx context.Context, expiredAccessToken *time.Time, metadata metadata.MD) (*BearerTokenResponse, error)
	ValidateAuthorizationRequest(ctx context.Context, decryptCode *AuthCodeToken, request *AuthCodeRequest) (*AuthorizationRequest, error)
	CompleteAuthorizationRequest(ctx context.Context, authRequest *AuthorizationRequest) (*string, *string, error)
	CreateAuthCode(ctx context.Context, expiredAuthCode *time.Time, client *entity.OAuthClient, userId *uuid.UUID, redirectUri *string, codeChallenge *string, codeChallengeMethod *string, scopes *string) (*entity.OAuthCode, error)
	RespondToRevokeRequest(ctx context.Context, grantType *string, token *string) error
}

// GenerateAuthorizationCodeRequest implements ServiceInterface.
func (g *ServiceServer) CompleteAuthorizationRequest(ctx context.Context, authRequest *AuthorizationRequest) (*string, *string, error) {

	if authRequest == nil {
		return nil, nil, fmt.Errorf("invalid authorization request")
	}

	if !authRequest.IsAuthApproved {
		return nil, nil, fmt.Errorf("auth is not approved")
	}

	if authRequest.User == nil {
		return nil, nil, fmt.Errorf("user is not found in request")
	}

	expired := g.dateInterval.GetEndDate()
	authCode, error := g.CreateAuthCode(ctx, &expired, authRequest.Client, &authRequest.User.Id, authRequest.RedirectUri, authRequest.CodeChallenge, (*string)(authRequest.CodeChallengeMethod), authRequest.Scopes)

	if error != nil {
		return nil, nil, error
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
	}

	authCodeToken, err := g.jwt.createAuthCodeToken(&payload)

	if err != nil {
		return nil, nil, err
	}

	redirectUri, err := g.makeRedirectUrl(authRequest.RedirectUri, map[string]interface{}{"code": authCodeToken})

	if err != nil {
		return nil, nil, err
	}

	return authCodeToken, redirectUri, nil
}

// CreateAuthCode implements ServiceInterface.
func (g *ServiceServer) CreateAuthCode(ctx context.Context, expiredAuthCode *time.Time, client *entity.OAuthClient, userId *uuid.UUID, redirectUri *string, codeChallenge *string, codeChallengeMethod *string, scopes *string) (*entity.OAuthCode, error) {
	splitString := strings.Split(*scopes, " ")
	var scoped = make([]string, len(splitString))
	for i, v := range splitString {
		v = strings.ToLower(v)
		scoped[i] = v
	}

	var user *entity.OAuthUser
	var errUser error
	if userId != nil {
		user, errUser = g.repository.GetUserByIdentifier(ctx, userId)
		if errUser != nil {
			return nil, errUser
		}
	}

	authCode := entity.OAuthCode{
		ExpairedAt:          *expiredAuthCode,
		RedirectUri:         redirectUri,
		CodeChallenge:       codeChallenge,
		CodeChallengeMethod: enum.CodeEnum(*codeChallengeMethod),
		Scopes:              scoped,
		User:                user,
		UserId:              userId,
	}

	err := g.repository.CreateOAuthCode(ctx, &authCode)

	if err != nil {
		return nil, err
	}

	return &authCode, nil
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
func (g *ServiceServer) ValidateAuthorizationRequest(ctx context.Context, decryptCode *AuthCodeToken, request *AuthCodeRequest) (*AuthorizationRequest, error) {
	panic("")
}
