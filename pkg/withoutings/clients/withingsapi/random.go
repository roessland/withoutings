package withingsapi

import (
	"crypto/rand"
	"encoding/base64"
	"regexp"
)

var b64nonAlphaNumeric = regexp.MustCompile(`[=+/]`)

func RandomNonce() string {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		panic(err)
	}

	return b64nonAlphaNumeric.ReplaceAllString(base64.RawURLEncoding.EncodeToString(b), "")
}
