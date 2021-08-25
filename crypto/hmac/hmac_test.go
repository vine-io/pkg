package hmac

import "testing"

func TestHmac(t *testing.T) {
	t.Log(Hmac("hello", "hmac"))
	t.Log(Hmac("hello", "hmac", "salt"))
	t.Log(HmacSha256("hello", "hmac"))
	t.Log(HmacSha256("hello", "hmac", "salt"))
}
