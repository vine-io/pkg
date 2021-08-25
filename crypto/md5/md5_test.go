package md5

import "testing"

func TestMD5(t *testing.T) {
	a := MD5("Hello md5")
	t.Log(a)

	a = MD5("Hello md5", "salt")
	t.Log(a)
}
