package service

import (
	"context"
	"fmt"
	"slices"
	"time"

	"github.com/gdochadipa/oauth2-go-project/internal/entity"
	"github.com/gdochadipa/oauth2-go-project/internal/enum"
	"github.com/gdochadipa/oauth2-go-project/internal/util"
	"github.com/gdochadipa/oauth2-go-project/pkg/api/v1/pb"
	"github.com/google/uuid"
)

type AuthorizationRequest struct {
	Scopes              []string
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
	UserID              *string        `json:"user_id"`
	RedirectURI         *string        `json:"redirect_uri"`
	CodeChallenge       *string        `json:"code_challenge"`
	CodeChallengeMethod *enum.CodeEnum `json:"code_challenge_method"`
	Audience            []string       `json:"audience"`
}

type AuthCodeGrantInterface interface {
	AccessTokenVerified(ctx context.Context, request *pb.AuthCodeGrantRequest, expiredAccessToken *time.Time) (*BearerTokenResponse, error)
	ValidateAuthorizationRequest(ctx context.Context, request *pb.GenerateCodeRequest) (*AuthorizationRequest, error)
	CompleteAuthorizationRequest(ctx context.Context, authRequest *AuthorizationRequest) (*string, *string, error)
	CreateAuthCode(ctx context.Context, expiredAuthCode *time.Time, client *entity.OAuthClient, userId *uuid.UUID, redirectUri *string, codeChallenge *string, codeChallengeMethod *string, scopes []string) (*entity.OAuthCode, error)
	RespondToRevokeRequest(ctx context.Context, grantType *string, token *string) error
	ValidateAuthorizationCode(ctx context.Context, redirectUri []string, client *entity.OAuthClient, authCodeToken *AuthCodeToken) error
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
	userId := authCode.User.Id.String()
	payload := PayloadAuthenticationCode{
		ClidenID:            authCode.Client.Id.String(),
		RedirectURI:         authCode.RedirectUri,
		AuthCodeID:          authCode.Code,
		Scopes:              authCode.Scopes,
		UserID:              &userId,
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
func (g *ServiceServer) CreateAuthCode(ctx context.Context, expiredAuthCode *time.Time, client *entity.OAuthClient, userId *uuid.UUID, redirectUri *string, codeChallenge *string, codeChallengeMethod *string, scopes []string) (*entity.OAuthCode, error) {
	// change from protobuf type string to []string
	// splitString := strings.Split(*scopes, " ")
	// var scoped = make([]string, len(splitString))
	// for i, v := range splitString {
	// 	v = strings.ToLower(v)
	// 	scoped[i] = v
	// }

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
		Scopes:              scopes,
		User:                user,
		UserId:              userId,
	}

	err := g.repository.CreateOAuthCode(ctx, &authCode)

	if err != nil {
		return nil, err
	}

	return &authCode, nil
}

// respondToAccessTokenRequest in oauth code
// need set metadata
// authorization => Bearer Base64(clientid:clientSecret)
func (g *ServiceServer) AccessTokenVerified(ctx context.Context, request *pb.AuthCodeGrantRequest, expiredAccessToken *time.Time) (*BearerTokenResponse, error) {
	client, redirectUri, errValidate := g.validateClient(ctx)

	if errValidate != nil {
		return nil, errValidate
	}

	// not need validation
	// it depend on gRPC proto
	// if required, it won't be null

	authCodeToken, errDecode := g.jwt.decodeAuthCodeToken(request.Code)

	if errDecode != nil {
		return nil, errDecode
	}

	err := g.ValidateAuthorizationCode(ctx, redirectUri, client, authCodeToken)

	if err != nil {
		return nil, err
	}

	var user *entity.OAuthUser = &entity.OAuthUser{}

	userID := uuid.MustParse(*authCodeToken.UserID)
	if authCodeToken.UserID != nil {
		user, err = g.repository.GetUserByIdentifier(ctx, &userID)
		if err != nil {
			return nil, err
		}
	}

	scopes, err := g.validateScopes(ctx, authCodeToken.Scopes)

	if err != nil {
		return nil, err
	}
	// need verfify scope with our scope
	userScope, errScope := g.repository.GetScopeByIdentifier(ctx, "authorization_code", scopes, client, &userID)

	if errScope != nil {
		return nil, errScope
	}

	authCode, err := g.repository.FindByCode(ctx, authCodeToken.AuthCodeID)

	if err != nil {
		return nil, err
	}

	if authCode.CodeChallenge != nil {
		if authCodeToken.CodeChallenge == nil {
			return nil, fmt.Errorf("empty.code.challenge")
		}

		if authCodeToken.CodeChallenge != authCode.CodeChallenge {
			return nil, fmt.Errorf("not.match.auth.code")
		}
		md, err := g.GettingMetadata(ctx)

		if err != nil {
			return nil, err
		}

		codeVerifier := md.Get("code_verifier")

		if len(codeVerifier) == 0 {
			return nil, fmt.Errorf("empty.code.verifier.metadata")
		}

		if !util.ValidateRegexCode(&codeVerifier[0]) {
			return nil, fmt.Errorf("invalid.code>.verifier")
		}

		var codeChallenge enum.CodeEnum = enum.Plain
		if authCodeToken.CodeChallengeMethod != nil {
			codeChallenge = *authCodeToken.CodeChallengeMethod
		}
		verifier, errVeri := util.NewVerifyCodeChallenge().MethodVerify((*string)(&codeChallenge), &codeVerifier[0], authCode.CodeChallenge)

		if errVeri != nil {
			return nil, errVeri
		}

		if !verifier {
			return nil, fmt.Errorf("failed.to.verify.code.challenge")
		}
	}

	// create access token and save on repository
	// !check if access token must generate by random or something better
	accessToken, errorIssue := g.issueAccessToken(ctx, client, user, userScope, expiredAccessToken)

	if errorIssue != nil {
		return nil, errorIssue
	}

	var errorRefresh error
	// create new refresh token
	// update access token
	// still generate it by rendom
	accessToken, errorRefresh = g.issueRefreshToken(ctx, accessToken, client)

	if errorRefresh != nil {
		return nil, errorRefresh
	}

	// revoke auth code
	// made it expiry, or update status is revoked
	g.repository.RevokedOAuthCode(ctx, authCodeToken.AuthCodeID)

	bearerTokenResponse, errBearer := g.makeBearerTokenResponse(ctx, client, accessToken, userScope, user)

	if errBearer != nil {
		return nil, errBearer
	}
	return bearerTokenResponse, nil
}

// RespondToRevokeRequest implements ServiceInterface.
func (g *ServiceServer) RespondToRevokeRequest(ctx context.Context, grantType *string, token *string) error {
	panic("unimplemented")
}

// ValidateAuthorizationRequest implements ServiceInterface.
func (g *ServiceServer) ValidateAuthorizationRequest(ctx context.Context, request *pb.GenerateCodeRequest) (*AuthorizationRequest, error) {
	// client id is required on protobuf, no need validation
	// because the protobuf set require the clientID
	if err := uuid.Validate(request.ClientId); err != nil {
		return nil, err
	}
	oauthClient, errorGet := g.repository.OAuthClientFindById(ctx, uuid.MustParse(request.ClientId))

	if errorGet != nil {
		return nil, errorGet
	}

	if oauthClient == nil {
		return nil, fmt.Errorf("oauthclient.is.not.found")
	}

	redirectUri, err := g.getRedirectUri(request, oauthClient)

	if err != nil {
		return nil, err
	}

	authRequest := AuthorizationRequest{
		Scopes:        request.Scopes,
		Client:        oauthClient,
		RedirectUri:   redirectUri,
		CodeChallenge: &request.CodeChallenge,
	}

	// for default we will require requiresPKCE
	// protobuf for code flow is required  PKCE

	codeChallengeMethod, err := enum.ConvertCodeEnum(&request.CodeChallengeMethod)

	if err != nil {
		return nil, err
	}

	if request.CodeChallengeMethod == "plain" {
		return nil, fmt.Errorf("code.challenge.method.must.S256")
	}

	authRequest.CodeChallengeMethod = &codeChallengeMethod

	return &authRequest, nil
}

func (g *ServiceServer) ValidateAuthorizationCode(ctx context.Context, redirectUri []string, client *entity.OAuthClient, authCodeToken *AuthCodeToken) error {
	if authCodeToken.AuthCodeID == "" {
		return fmt.Errorf("auth code id are undefined")
	}

	isAuthCodeRevoked, err := g.repository.IsRevokedOAuthCode(ctx, authCodeToken.AuthCodeID)

	if err != nil {
		return err
	}

	now := time.Now().UTC().Unix()

	if now > authCodeToken.ExpiresAt.Unix() || isAuthCodeRevoked {
		return fmt.Errorf("Authorization code is expired or revoked")
	}

	if *client.ClientId != uuid.MustParse(authCodeToken.ClientID) {
		return fmt.Errorf("miss match client id")
	}

	if authCodeToken.RedirectURI == nil {
		return fmt.Errorf("invalid.token.redirect.uri")
	}

	isMatch := slices.Contains(redirectUri, *authCodeToken.RedirectURI)
	if !isMatch {
		return fmt.Errorf("redirect.uri.unmatch")
	}

	return nil
}
