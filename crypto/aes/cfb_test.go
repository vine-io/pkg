package aes

import "testing"

func TestCFB(t *testing.T) {
	source := "hello world"

	out, err := CFBEncode(source)
	if err != nil {
		t.Fatal(err)
	}

	orig, err := CFBDecode(out)
	if err != nil {
		t.Fatal(err)
	}

	if source != orig {
		t.Log()
	}

	t.Log(out)
}
