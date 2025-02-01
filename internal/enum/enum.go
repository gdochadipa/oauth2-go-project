package enum

type CodeEnum string

const (
	S256  CodeEnum = "S256"
	Plain CodeEnum = "plain"
)

type GrantIdentifier string

const (
	Code          GrantIdentifier = "authorization_code"
	Client        GrantIdentifier = "client_credentials"
	Refresh_token GrantIdentifier = "refresh_token"
	Password      GrantIdentifier = "password"
	Implicit      GrantIdentifier = "implicit"
)

var grantIdentifierType = map[GrantIdentifier]struct{}{
	Code:          {},
	Client:        {},
	Refresh_token: {},
	Password:      {},
	Implicit:      {},
}

func IsIncludeGrant(s GrantIdentifier) bool {
	_, isExists := grantIdentifierType[s]
	return isExists
}
