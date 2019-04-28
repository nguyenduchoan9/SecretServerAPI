package repo

import (
	"crypto/sha1"
	"encoding/base64"
)

func GenerateHashBy(text string) string {
	hash := sha1.New()
	hash.Write([]byte(text))
	return base64.URLEncoding.EncodeToString(hash.Sum(nil))
}
