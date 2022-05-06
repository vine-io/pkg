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
