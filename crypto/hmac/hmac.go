package hmac

import (
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/hex"
	"hash"
)

func Hmac(key string, data string, salt ...string) string {
	return mac(md5.New, key, data, salt...)
}

func HmacSha256(key, data string, salt ...string) string {
	return mac(sha256.New, key, data, salt...)
}

func HmacSha512(key, data string, salt ...string) string {
	return mac(sha512.New, key, data, salt...)
}

func mac(fn func() hash.Hash, key string, data string, salt ...string) string {
	h := hmac.New(fn, []byte(key))
	h.Write([]byte(data))
	b := ""
	if salt != nil {
		b = salt[0]
	}
	return hex.EncodeToString(h.Sum([]byte(b)))
}
