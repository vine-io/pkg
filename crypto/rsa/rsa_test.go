package rsa

import (
	"encoding/hex"
	"testing"
)

func TestEncode(t *testing.T) {
	source := []byte("hello world")

	out, err := Encode([]byte(source))
	if err != nil {
		t.Fatal(err)
	}

	orig, err := Decode(out)
	if err != nil {
		t.Fatal(err)
	}

	if string(source) != string(orig) {
		t.Log()
	}

	t.Log(hex.EncodeToString(out))
}
