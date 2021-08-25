package aes

import (
	"testing"
)

func TestCRT(t *testing.T) {
	source := "hello world"

	out, err := CRT(source)
	if err != nil {
		t.Fatal(err)
	}

	orig, err := CRT(out)
	if err != nil {
		t.Fatal(err)
	}

	if source != orig {
		t.Log()
	}

	t.Log(out)
}
