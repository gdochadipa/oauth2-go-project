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

// encryptAccessToken implements GrantInterface.
func (g *ServiceServer) encryptAccessToken() {
	panic("unimplemented")
}

// encryptRefreshToken implements GrantInterface.
func (g *ServiceServer) encryptRefreshToken() {
	panic("unimplemented")
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
