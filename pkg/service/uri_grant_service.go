package service

import (
	"github.com/gdochadipa/oauth2-go-project/internal/entity"
)

type UriGrantService interface {
	getRedirectUri(redirectUri string, client *entity.OAuthClient) (*string, error)
	validateRedirectUri(redirectUri string, client *entity.OAuthClient) (*string, error)
	makeRedirectUrl(uri *string, params map[string]any) (*string, error)
}

// getRedirectUri implements ServiceInterface.
func (g *ServiceServer) getRedirectUri(redirectUri string, client *entity.OAuthClient) (*string, error) {
	panic("unimplemented")
}

// validateRedirectUri implements ServiceInterface.
func (g *ServiceServer) validateRedirectUri(redirectUri string, client *entity.OAuthClient) (*string, error) {
	panic("unimplemented")
}

func (g *ServiceServer) makeRedirectUrl(uri *string, params map[string]any) (*string, error) {
	panic("unimplemented")
}
