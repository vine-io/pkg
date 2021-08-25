package md5

import (
	"crypto/md5"
	"encoding/hex"
)

func MD5(s string, salt ...string) string {
	h := md5.New()
	h.Write([]byte(s))
	b := ""
	if salt != nil {
		b = salt[0]
	}
	return hex.EncodeToString(h.Sum([]byte(b)))
}
