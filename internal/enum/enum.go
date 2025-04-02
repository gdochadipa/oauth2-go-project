package enum

import (
	"fmt"
)

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

func ConvertCodeEnum(s *string) (CodeEnum, error) {
	switch *s {
	case "S256":
		return S256, nil
	case "Plain":
		return Plain, nil

	default:
		return Plain, fmt.Errorf("not.found.code.enum")
	}
}
