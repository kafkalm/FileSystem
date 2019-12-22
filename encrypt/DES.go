package encrypt

import (
	"bytes"
	"crypto/cipher"
	"crypto/des"
)

func DESEncrypt(plaintext,key []byte) ([]byte,error) {
	block,err := des.NewCipher(key)
	if err != nil {
		return nil,err
	}
	plaintext = PKCS5Padding(plaintext,block.BlockSize())
	blockmode := cipher.NewCBCEncrypter(block,key)
	ciphertext := make([]byte,len(plaintext))
	blockmode.CryptBlocks(ciphertext,plaintext)
	return ciphertext,nil
}

func DESDecrypt(ciphertext,key []byte) ([]byte,error) {
	block,err := des.NewCipher(key)
	if err != nil {
		return nil,err
	}
	blockmode := cipher.NewCBCDecrypter(block,key)
	plaintext := make([]byte,len(ciphertext))
	blockmode.CryptBlocks(plaintext,ciphertext)
	plaintext = PKCS5UnPadding(plaintext)
	return plaintext,nil
}

func PKCS5Padding(cipherText []byte, blockSize int) []byte {
	padding := blockSize - len(cipherText) % blockSize
	padText := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(cipherText, padText...)
}

func PKCS5UnPadding(origData []byte) []byte {
	length := len(origData)
	// 去掉最后一个字节 unpadding 次
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}