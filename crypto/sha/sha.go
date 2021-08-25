package sha

import (
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/hex"
	"hash"
)

func Sha1(s string, salt ...string) string {
	return sha(sha1.New, s, salt...)
}

func Sha256(s string, salt ...string) string {
	return sha(sha256.New, s, salt...)
}

func Sha512(s string, salt ...string) string {
	return sha(sha512.New, s, salt...)
}

func sha(fn func() hash.Hash, s string, salt ...string) string {
	h := fn()
	h.Write([]byte(s))
	b := ""
	if salt != nil {
		b = salt[0]
	}
	return hex.EncodeToString(h.Sum([]byte(b)))
}
