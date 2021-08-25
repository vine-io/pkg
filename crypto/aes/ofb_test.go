package aes

import "testing"

func TestOFB(t *testing.T) {
	source := "hello world"

	out, err := OFBEncode(source)
	if err != nil {
		t.Fatal(err)
	}

	orig, err := OFBDecode(out)
	if err != nil {
		t.Fatal(err)
	}

	if source != orig {
		t.Log()
	}

	t.Log(out)
}
