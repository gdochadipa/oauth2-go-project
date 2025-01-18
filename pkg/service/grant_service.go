package service

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

type GrantService struct {
	jwt JWTInterface
}

// encryptAccessToken implements GrantInterface.
func (g *GrantService) encryptAccessToken() {
	panic("unimplemented")
}

// encryptRefreshToken implements GrantInterface.
func (g *GrantService) encryptRefreshToken() {
	panic("unimplemented")
}

// getBasicAuthCredentials implements GrantInterface.
func (g *GrantService) getBasicAuthCredentials() {
	panic("unimplemented")
}

// getClientCredentials implements GrantInterface.
func (g *GrantService) getClientCredentials() {
	panic("unimplemented")
}

// getGrantType implements GrantInterface.
func (g *GrantService) getGrantType() {
	panic("unimplemented")
}

// issueAccessToken implements GrantInterface.
func (g *GrantService) issueAccessToken() {
	panic("unimplemented")
}

// makeBearerTokenResponse implements GrantInterface.
func (g *GrantService) makeBearerTokenResponse() {
	panic("unimplemented")
}

// validateClient implements GrantInterface.
func (g *GrantService) validateClient() {
	panic("unimplemented")
}

// validateScopes implements GrantInterface.
func (g *GrantService) validateScopes() {
	panic("unimplemented")
}

func NewGrantRepository(jwt JWTInterface) GrantInterface {
	return &GrantService{jwt: jwt}
}
