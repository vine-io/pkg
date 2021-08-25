package aes

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"errors"
	"io"
)

var (
	DefaultCFBKey = "ABCDEFGHIJKLMNO1"
)

func CFBEncode(text string) (string, error) {
	block, err := aes.NewCipher([]byte(DefaultCFBKey))
	if err != nil {
		return "", err
	}
	out := make([]byte, aes.BlockSize+len(text))
	iv := out[:aes.BlockSize]
	if _, err = io.ReadFull(rand.Reader, iv); err != nil {
		return "", err
	}
	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(out[aes.BlockSize:], []byte(text))
	return string(out), nil
}

func CFBDecode(text string) (string, error) {
	block, _ := aes.NewCipher([]byte(DefaultCFBKey))
	if len(text) < aes.BlockSize {
		return "", errors.New("ciphertext too short")
	}
	iv := text[:aes.BlockSize]
	text = text[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, []byte(iv))
	b := []byte(text)
	stream.XORKeyStream(b, b)
	return string(b), nil
}
