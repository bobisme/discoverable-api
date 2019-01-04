package models

import (
	"crypto/rand"
	"encoding/base64"
)

func newId() string {
	out := make([]byte, 8)
	rand.Read(out)
	return base64.RawURLEncoding.EncodeToString(out)
}
