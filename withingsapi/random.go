package withingsapi

import (
	"crypto/rand"
	"encoding/base64"
	"strings"
)

func RandomNonce() string {
	b := make([]byte, 9)
	_, err := rand.Read(b)
	if err != nil {
		panic(err)
	}

	return strings.TrimRight(base64.StdEncoding.EncodeToString(b), "=")
}
