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
