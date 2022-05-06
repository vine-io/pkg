// MIT License
//
// Copyright (c) 2021 Lack
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

package rsa

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
)

var (
	// DefaultPrivateKey private key
	// openssl genrsa -out keys.pem 1024
	DefaultPrivateKey = []byte(`
-----BEGIN RSA PRIVATE KEY-----
MIICXQIBAAKBgQDkv765UFFGyCJJWlPT+2ZeOxsw9ZPmQGniwq0/blxwvRGUlfzo
C4x8t2QY0crm1sZP6P93HXjg09bOsHbQT04itoyCY977/dd0nBIgZ/qVkehymMoA
tclkoAqNzTXYNzp5DZ6kRMtlmyX7EgndRPQ3Mm9cNd5q+paOUCAwEiddGwIDAQAB
AoGAKo7mDBI+XN3lSyJsEDdor0Vt5Kj78E2xpYe0teVxe2QhvjQ9jmp/o8B29gsq
JdJ1qO5fgSjRkXq4L1IzeMQYdBDMxqG9CGSufWll0LtSmNAIBm6AKNO4dA74OzpC
UO7nzX+djGb02ZG5tKRQ4mMuLW/2PwoSepfWccwAzc7np1ECQQD5/v9bUqtaz7Cw
eyMBLpNNp9sJNS0RTfz9EwpIyynOV8CvOJzvRfHGb2xtHGqSIFD2ptb5zysBe9/v
D46HTAIpAkEA6j4eQ7Ms2GH7TyV8EL/0WrM39OTa85Z5DdmBvkpSrM/mTGd/e0mF
E8c9tJ8JFswTdIqKj5HEEjF4GJNKesO1owJBAIDMnPmLBRe7a3fxaR6BxYi704DR
8c85k/87IRBSA874rSBZlZk9OwyWeZFZk5qHpc7+NEHuN2UDUmNTa4ZPZckCQQCn
UqQPvAeGscbwbFhJJrUHrQmFl4yHf68NI5e4NCMGaqOZZDz99jBnRmVfhlLZxAEJ
uITttTQXwtqEw4HqW659AkBNZIelmCJL9zFV1VcXOgzuO870a2zm/hodxs9ocndk
2BENmtxu78U6IHLB3GzuWUiBXP1RLms/4Vd3Q4MxUyyb
-----END RSA PRIVATE KEY-----
`)

	// DefaultPublicKey public key
	// openssl rsa -in keys.pem -pubout -out public.pem
	DefaultPublicKey = []byte(`
-----BEGIN PUBLIC KEY-----
MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQDkv765UFFGyCJJWlPT+2ZeOxsw
9ZPmQGniwq0/blxwvRGUlfzoC4x8t2QY0crm1sZP6P93HXjg09bOsHbQT04itoyC
Y977/dd0nBIgZ/qVkehymMoAtclkoAqNzTXYNzp5DZ6kRMtlmyX7EgndRPQ3Mm9c
Nd5q+paOUCAwEiddGwIDAQAB
-----END PUBLIC KEY-----
`)
)

func Encode(text []byte) ([]byte, error) {
	// 解密 pem 格式的公钥
	block, _ := pem.Decode(DefaultPublicKey)
	if block == nil {
		return nil, errors.New("public key error")
	}
	// 解析公钥
	pubInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	// 类型断言
	pub := pubInterface.(*rsa.PublicKey)
	// 加密
	return rsa.EncryptPKCS1v15(rand.Reader, pub, text)
}

func Decode(text []byte) ([]byte, error) {
	// 解密
	block, _ := pem.Decode(DefaultPrivateKey)
	if block == nil {
		return nil, errors.New("private key error!")
	}
	// 解析 PKCS1 格式的私钥
	priv, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	// 解密
	return rsa.DecryptPKCS1v15(rand.Reader, priv, text)
}
