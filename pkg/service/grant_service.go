package service

import (
	"fmt"
	"time"

	"github.com/gdochadipa/oauth2-go-project/internal/entity"
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

// getBasicAuthCredentials implements GrantInterface.
func (g *ServiceServer) getBasicAuthCredentials() {
	panic("unimplemented")
}

// getClientCredentials implements GrantInterface.
func (g *ServiceServer) getClientCredentials() {
	panic("unimplemented")
}

// getGrantType implements GrantInterface.
func (g *ServiceServer) getGrantType() {
	panic("unimplemented")
}

// issueAccessToken implements GrantInterface.
func (g *ServiceServer) issueAccessToken() {
	panic("unimplemented")
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
func (g *ServiceServer) validateScopes() {
	panic("unimplemented")
}
