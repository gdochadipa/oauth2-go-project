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
	"google.golang.org/grpc/metadata"
)

type GrantInterface interface {
	makeBearerTokenResponse()
	encryptRefreshToken()
	encryptAccessToken()
	validateClient()
	getClientCredentials()
	getBasicAuthCredentials()
	validateScopes()
	issueAccessToken()
	getGrantType()
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
func (g *ServiceServer) encryptRefreshToken(client *entity.OAuthClient, refreshToken *entity.OAuthToken, scopes []entity.OAuthScope) (*string, error) {
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
*/
func (g *ServiceServer) getBasicAuthCredentials(metadata metadata.MD) (*string, *string) {
	if metadata.Get("authorization") == nil {
		return nil, nil
	}

	auth := metadata.Get("authorization")[0]
	if strings.HasPrefix(auth, "Basic ") || auth == "" {
		return nil, nil
	}

	decoded := util.Base64Decode(auth[6:])
	if !strings.Contains(decoded, ":") {
		return nil, nil
	}

	result := strings.Split(decoded, ":")
	if result[0] != "" && result[1] != "" {
		return &result[0], &result[1]
	}

	return nil, nil
}

// getClientCredentials implements GrantInterface.
func (g *ServiceServer) getClientCredentials(ctx context.Context, metadata metadata.MD) (string, string, error) {
	// basicAuthUser, basicAuthPass := g.getBasicAuthCredentials(metadata)
	panic("unimplemented")
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
func (g *ServiceServer) issueAccessToken(ctx context.Context, client *entity.OAuthClient, user *entity.OAuthUser, scopes []entity.OAuthScope) (*entity.OAuthToken, error) {
	accessToken, error := g.repository.IssueToken(ctx, client, scopes, user)

	if error != nil {
		return nil, error
	}
	accessToken.AccessTokenExpiresAt = time.Now().Add(time.Hour * 24)
	g.repository.PersistAccessToken(ctx, accessToken)

	return accessToken, nil
}

// makeBearerTokenResponse implements GrantInterface.
func (g *ServiceServer) makeBearerTokenResponse() {
	panic("unimplemented")
}

// validateClient implements GrantInterface.
func (g *ServiceServer) validateClient() {
	panic("unimplemented")
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
