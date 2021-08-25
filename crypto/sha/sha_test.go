package sha

import (
	"testing"
)

func TestSha1(t *testing.T) {
	t.Log(Sha1("1"))
}

func TestSha256(t *testing.T) {
	t.Log(Sha256("1"))
}

func TestSha512(t *testing.T) {
	t.Log(Sha512("1"))
}
