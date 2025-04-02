package service

import (
	"context"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/gdochadipa/oauth2-go-project/internal/entity"
	"github.com/gdochadipa/oauth2-go-project/internal/enum"
	"github.com/gdochadipa/oauth2-go-project/internal/util"
	"github.com/google/uuid"
	"google.golang.org/grpc/metadata"
)

type GrantInterface interface {
	makeBearerTokenResponse(ctx context.Context, client *entity.OAuthClient, accessToken *entity.OAuthToken, scopes []entity.OAuthScope, user *entity.OAuthUser) (*BearerTokenResponse, error)
	encryptRefreshToken(ctx context.Context, client *entity.OAuthClient, refreshToken *entity.OAuthToken, scopes []entity.OAuthScope) (*string, error)
	encryptAccessToken(username string, client *entity.OAuthClient, accessToken *entity.OAuthToken, scopes []entity.OAuthScope) (*string, error)
	validateClient(ctx context.Context) (*entity.OAuthClient, []string, error)
	getClientCredentials(ctx context.Context, metadata metadata.MD) (*string, *string, error)
	getBasicAuthCredentials(metadata metadata.MD) (*string, *string, error)
	validateScopes(ctx context.Context, scopes []string) ([]entity.OAuthScope, error)
	issueAccessToken(ctx context.Context, client *entity.OAuthClient, user *entity.OAuthUser, scopes []entity.OAuthScope, expiredAccessToken *time.Time) (*entity.OAuthToken, error)
	getGrantType(ctx context.Context, metadata metadata.MD) (*enum.GrantIdentifier, error)
	issueRefreshToken(ctx context.Context, accessToken *entity.OAuthToken, client *entity.OAuthClient) (*entity.OAuthToken, error)
	GettingMetadata(ctx context.Context) (*metadata.MD, error)
}

type BearerTokenResponse struct {
	TokenType          string
	ExpiresIn          int16
	AccessToken        *string
	RefreshAccessToken *string
	Scopes             []entity.OAuthScope
}

// encryptAccessToken implements GrantInterface.
func (g *ServiceServer) encryptAccessToken(username string, client *entity.OAuthClient, accessToken *entity.OAuthToken, scopes []entity.OAuthScope) (*string, error) {

	stringScope := func() string {
		if len(scopes) == 0 {
			return ""
		}

		result := scopes[0].Id
		for _, scope := range scopes[1:] {
			result += scope.Id
		}
		return result
	}()

	now := time.Now()

	if accessToken.UserId == nil {
		return nil, fmt.Errorf("user id is not found.")
	}

	return g.jwt.createAccessToken(*accessToken.UserId, username, stringScope, accessToken.ClientId.String(), accessToken.AccessTokenExpiresAt, client.RedirectUris[0], now, now.Add(1*time.Hour))

}

// encryptRefreshToken implements GrantInterface.
func (g *ServiceServer) encryptRefreshToken(ctx context.Context, client *entity.OAuthClient, refreshToken *entity.OAuthToken, scopes []entity.OAuthScope) (*string, error) {
	stringScope := func() string {
		if len(scopes) == 0 {
			return ""
		}

		result := scopes[0].Id
		for _, scope := range scopes[1:] {
			result += scope.Id
		}
		return result
	}()

	if refreshToken.UserId == nil {
		return nil, fmt.Errorf("user id is not found.")
	}

	return g.jwt.createRefreshToken(*refreshToken.UserId, client.ClientId.String(), stringScope, *refreshToken.RefreshTokenExpiresAt)
}

/*
*

	this function for getting basic credentials from metadata
	because get authorization from metadata
*/
func (g *ServiceServer) getBasicAuthCredentials(metadata metadata.MD) (*string, *string, error) {
	if metadata.Get("authorization") == nil {
		return nil, nil, fmt.Errorf("invalid.credentials.format")
	}

	auth := metadata.Get("authorization")[0]
	if strings.HasPrefix(auth, "Basic ") || auth == "" {
		return nil, nil, fmt.Errorf("invalid.credentials.format")
	}

	decoded := util.Base64Decode(auth[6:])
	if !strings.Contains(decoded, ":") {
		return nil, nil, fmt.Errorf("invalid.credentials.format")
	}

	result := strings.Split(decoded, ":")
	if result[0] != "" && result[1] != "" {
		return &result[0], &result[1], fmt.Errorf("invalid.credentials.format")
	}

	return nil, nil, fmt.Errorf("invalid.credentials.format")
}

// getClientCredentials implements GrantInterface.
/**
Basic dan Bearer
untuk client credentials akan di taruh di  Bearer Token (metadata)
dengan di encode base64 format "client_id:client_secret"

alurnya bakal dimiripin sama aws cognito

kalo client credentials dan  itu pake Bearer Base64(clientid:clientSecret)

kalo pake Basic Base64(clientid:clientSecret)

ga mungkin bakal naruh credentials di metadata

- Opsi kedua adalah menggunakan handshake

tapi nanti aja opsi handshakenya, ribet bgt wkwk
*/
func (g *ServiceServer) getClientCredentials(ctx context.Context, metadata metadata.MD) (*string, *string, error) {
	basicAuthUser, basicAuthPass, error := g.getBasicAuthCredentials(metadata)

	if error != nil {
		return nil, nil, error
	}

	if basicAuthUser != nil || basicAuthPass != nil {
		return nil, nil, fmt.Errorf("credentials.is.empty")
	}

	return basicAuthUser, basicAuthPass, nil
}

/*
*
 */
func (g *ServiceServer) getGrantType(ctx context.Context, metadata metadata.MD) (*enum.GrantIdentifier, error) {
	if metadata.Get("grant_type") == nil {
		return nil, fmt.Errorf("grant.type.not.found")
	}

	grantType := metadata.Get("grant_type")[0]

	if !enum.IsIncludeGrant(enum.GrantIdentifier(grantType)) {
		return nil, fmt.Errorf("grant.not.include.enum")
	}

	return (*enum.GrantIdentifier)(&grantType), nil
}

// issueAccessToken implements GrantInterface.
func (g *ServiceServer) issueAccessToken(ctx context.Context, client *entity.OAuthClient, user *entity.OAuthUser, scopes []entity.OAuthScope, expiredAccessToken *time.Time) (*entity.OAuthToken, error) {
	// create access token
	//CreateDataAccessToken
	accessToken, error := g.repository.CreateDataAccessToken(ctx, client, scopes, user)

	if error != nil {
		return nil, error
	}
	// accessToken.AccessTokenExpiresAt = time.Now().Add(time.Hour * 24)
	accessToken.AccessTokenExpiresAt = *expiredAccessToken
	g.repository.PersistAccessToken(ctx, accessToken)

	return accessToken, nil
}

func (g *ServiceServer) issueRefreshToken(ctx context.Context, accessToken *entity.OAuthToken, client *entity.OAuthClient) (*entity.OAuthToken, error) {
	//IssueRefreshToken
	// update refresh token
	return g.repository.UpdateRefreshToken(ctx, accessToken, client)
}

// makeBearerTokenResponse implements GrantInterface.
func (g *ServiceServer) makeBearerTokenResponse(ctx context.Context, client *entity.OAuthClient, accessToken *entity.OAuthToken, scopes []entity.OAuthScope, user *entity.OAuthUser) (*BearerTokenResponse, error) {
	encryptedAccessToken, error := g.encryptAccessToken(user.Username, client, accessToken, scopes)

	if error != nil {
		return nil, error
	}

	var encryptedRefreshToken *string
	if accessToken.RefreshToken != nil {
		encryptedRefreshToken, error = g.encryptRefreshToken(ctx, client, accessToken, scopes)

		if error != nil {
			return nil, error
		}
	}

	return &BearerTokenResponse{
		TokenType:          "Bearer",
		ExpiresIn:          int16(accessToken.AccessTokenExpiresAt.Unix()),
		AccessToken:        encryptedAccessToken,
		RefreshAccessToken: encryptedRefreshToken,
		Scopes:             scopes,
	}, nil
}

// validateClient implements GrantInterface.
// must in on bearer token
// header is metadata
func (g *ServiceServer) validateClient(ctx context.Context) (*entity.OAuthClient, []string,  error) {
	md, ok := metadata.FromIncomingContext(ctx)

	if !ok {
		return nil, nil,  fmt.Errorf("something wrong with metadata.")
	}

	clientId, clientSecret, errCredent := g.getClientCredentials(ctx, md)

	if errCredent != nil {
		return nil, nil,  errCredent
	}

	grantType, getError := g.getGrantType(ctx, md)

	if getError != nil {
		return nil, nil, getError
	}

	client, repoError := g.repository.OAuthClientFindById(ctx, uuid.MustParse(*clientId))

	if repoError != nil {
		return nil,nil,  repoError
	}

	if clientSecret != nil {
		return nil,nil, fmt.Errorf("invalid.credentials")
	}

	userValidationSuccess := g.repository.IsClientValid(ctx, grantType, client, clientSecret)

	if !userValidationSuccess {
		return nil,nil, fmt.Errorf("invalid.credentials.not.pass")
	}

	redirectUri := md.Get("redirectUri");
	if md.Get("redirectUri") == nil {
		return nil,nil, fmt.Errorf("missing.redirect.uri")
	}

	return client, redirectUri,  nil

}

// validateScopes implements GrantInterface.
func (g *ServiceServer) validateScopes(ctx context.Context, scopes []string) ([]entity.OAuthScope, error) {

	validScopes, error := g.repository.GetAllScope(ctx, scopes)

	mapingScope := make([]string, len(validScopes))

	for _, scope := range validScopes {
		mapingScope = append(mapingScope, scope.Name)
	}

	if reflect.DeepEqual(mapingScope, scopes) {
		return nil, fmt.Errorf("invalid.scope")
	}

	return validScopes, error
}

func (g *ServiceServer) GettingMetadata(ctx context.Context) (*metadata.MD, error) {
	md, ok := metadata.FromIncomingContext(ctx)

	if !ok {
		return nil, fmt.Errorf("something wrong with metadata.")
	}

	return &md, nil
}
