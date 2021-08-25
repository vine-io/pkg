package aes

import "testing"

func TestCBC(t *testing.T) {
	orig := "source"
	out := CBCEncode(orig)
	dstring := CBCDecode(out)

	if orig != dstring {
		t.Fatal("invalid cbc")
	}
	t.Log(out)
}
