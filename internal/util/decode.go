package util

import (
	"encoding/base64"
	"strings"
)

func Base64Decode(str string) string {
	sDec, _ := base64.StdEncoding.DecodeString(str)

	return string(sDec)
}

func Base64Encode(str string) string {
	sEnc := base64.StdEncoding.EncodeToString([]byte(str))
	return sEnc
}

func Base64UrlEncode(str string) string {
	encode := strings.NewReplacer("+", "-", "/", "_", "=", "").Replace(Base64Encode(str))
	return encode

}
