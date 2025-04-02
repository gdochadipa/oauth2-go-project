package util

import (
	"fmt"
)

type VerifyCodeChallenge interface {
	VerifyS256(codeVerify *string, codeChallenge *string) bool
	VerifyPlain(codeVerify *string, codeChallenge *string) bool
	MethodVerify(method *string, codeVerify *string, codeChallenge *string) (bool, error)
}

type VerifyCode struct {
}

// MethodVerify implements VerifyCodeChallenge.
func (v *VerifyCode) MethodVerify(method *string, codeVerify *string, codeChallenge *string) (bool, error) {
	switch *method {
	case "S256":
		return v.VerifyS256(codeVerify, codeChallenge), nil
	case "plain":
		return v.VerifyPlain(codeVerify, codeChallenge), nil
	default:
		return false, fmt.Errorf("empty.code.verifier.metadata")
	}

}

// VerifyPlain implements VerifyCodeChallenge.
func (v *VerifyCode) VerifyPlain(codeVerify *string, codeChallenge *string) bool {
	panic("unimplemented")
}

// VerifyS256 implements VerifyCodeChallenge.
func (v *VerifyCode) VerifyS256(codeVerify *string, codeChallenge *string) bool {
	panic("unimplemented")
}

func NewVerifyCodeChallenge() VerifyCodeChallenge {
	return &VerifyCode{}
}
