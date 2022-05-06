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
	"bytes"
	"crypto/aes"
	"crypto/cipher"
)

var (
	DefaultCRTKey = "1443flfsaWfdasds"
)

func CRT(text string) (string, error) {
	// 创建 cipher.Block 接口
	block, err := aes.NewCipher([]byte(DefaultCRTKey))
	if err != nil {
		return "", err
	}

	// 创建分组模式，在crypto/cipher包中
	iv := bytes.Repeat([]byte("a"), block.BlockSize())
	stream := cipher.NewCTR(block, iv)

	// 加密
	dst := make([]byte, len(text))
	stream.XORKeyStream(dst, []byte(text))

	return string(dst), nil
}
