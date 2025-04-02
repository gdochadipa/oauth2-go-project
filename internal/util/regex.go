package util

import "regexp"

func ValidateRegexCode(codeVerifier *string) bool {
	re := regexp.MustCompile(`/^[A-Za-z0-9-._~]{43,128}$/`)
	return re.MatchString(*codeVerifier)
}
