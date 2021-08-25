package aes

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"fmt"
	"io"
)

var (
	DefaultOFBKey = "ABCDEFGHIJKLMNO1" // 16bit or 32bit
)

func OFBEncode(text string) (string, error) {
	data := PKCS7Padding([]byte(text), aes.BlockSize)
	block, _ := aes.NewCipher([]byte(DefaultOFBKey))
	out := make([]byte, aes.BlockSize+len(data))
	iv := out[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", err
	}

	stream := cipher.NewOFB(block, iv)
	stream.XORKeyStream(out[aes.BlockSize:], data)
	return string(out), nil
}

func OFBDecode(text string) (string, error) {
	block, _ := aes.NewCipher([]byte(DefaultOFBKey))
	iv := text[:aes.BlockSize]
	data := text[aes.BlockSize:]
	if len(data)%aes.BlockSize != 0 {
		return "", fmt.Errorf("data is not a multiple of the block size")
	}

	out := make([]byte, len(data))
	mode := cipher.NewOFB(block, []byte(iv))
	mode.XORKeyStream(out, []byte(data))

	out = PKCS7UnPadding(out)
	return string(out), nil
}
