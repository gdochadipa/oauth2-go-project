package repository

import (
	"crypto/sha256"
	"encoding/base64"
)

type CodeVerifRepository struct {
	codeVerify    string
	codeChallenge string
}

func (r *CodeVerifRepository) S256Verifier() *bool {
	hash := sha256.Sum256([]byte(r.codeVerify))

	encode := base64.RawURLEncoding.EncodeToString(hash[:])

	verif := encode == r.codeChallenge

	return &verif
}

func (r *CodeVerifRepository) PlainVerify() *bool {
	verif := r.codeVerify == r.codeChallenge

	return &verif
}
