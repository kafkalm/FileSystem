package encrypt

import (
	"crypto/rsa"
	"crypto/rand"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"io/ioutil"
	"os"
)
// 生成密钥
func GenKey(bit int) (*rsa.PublicKey,*rsa.PrivateKey,error) {
	privatekey,err := rsa.GenerateKey(rand.Reader,bit)
	if err != nil {
		return nil,nil,err
	}
	publickey := &privatekey.PublicKey
	return publickey,privatekey,nil
}

// key保存到文件中
func DumpPrivateKeyToFile (privatekey *rsa.PrivateKey,filename string) (error) {
	var keybytes []byte = x509.MarshalPKCS1PrivateKey(privatekey)
	block := &pem.Block{
		Type : "RSA PRIVATE KEY",
		Bytes : keybytes,
	}
	file,err := os.Create(filename)
	if err != nil {
		return err
	}
	err = pem.Encode(file,block)
	if err != nil {
		return err
	}
	return nil
}

func DumpPublicKeyToFile (publickey *rsa.PublicKey,filename string) (error) {
	var keybytes []byte = x509.MarshalPKCS1PublicKey(publickey)
	block := &pem.Block{
		Type : "RSA PUBLIC KEY",
		Bytes : keybytes,
	}
	file,err := os.Create(filename)
	if err != nil {
		return err
	}
	err = pem.Encode(file,block)
	if err != nil {
		return err
	}
	return nil
}

// 从文件中读取key
func LoadPrivateKeyFromFile (filename string) (*rsa.PrivateKey,error) {
	buf,err := ioutil.ReadFile(filename)
	if err != nil {
		return nil,err
	}
	block,_ := pem.Decode(buf)
	if block == nil {
		return nil,errors.New("Load PrivateKey Failed")
	}
	privatekey,err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil,err
	}
	return privatekey,nil
}

func LoadPublicKeyFromFile (filename string) (*rsa.PublicKey,error) {
	buf,err := ioutil.ReadFile(filename)
	if err != nil {
		return nil,err
	}
	block,_ := pem.Decode(buf)
	if block == nil {
		return nil,errors.New("Load PublicKey Failed")
	}
	publickey,err := x509.ParsePKCS1PublicKey(block.Bytes)
	if err != nil {
		return nil,err
	}
	return publickey,nil
}

// key输出为base64编码字符串
func DumpPrivateKeyToBase64 (privatekey *rsa.PrivateKey) (string) {
	var keybytes []byte = x509.MarshalPKCS1PrivateKey(privatekey)

	keybase64 := base64.StdEncoding.EncodeToString(keybytes)
	return keybase64
}

func DumpPublicKeyToBase64 (publickey *rsa.PublicKey) (string) {
	var keybytes []byte = x509.MarshalPKCS1PublicKey(publickey)

	keybase64 := base64.StdEncoding.EncodeToString(keybytes)
	return keybase64
}

// base64编码字符串解码为key
func LoadPrivateKeyFromBase64 (keybase64 string) (*rsa.PrivateKey,error) {
	keybytes,err := base64.StdEncoding.DecodeString(keybase64)
	if err != nil {
		return nil,err
	}
	privatekey,err := x509.ParsePKCS1PrivateKey(keybytes)
	if err != nil {
		return nil,err
	}
	return privatekey,nil
}

func LoadPublicKeyFromBase64 (keybase64 string) (*rsa.PublicKey,error) {
	keybytes,err := base64.StdEncoding.DecodeString(keybase64)
	if err != nil {
		return nil,err
	}
	publickey,err := x509.ParsePKCS1PublicKey(keybytes)
	if err != nil {
		return nil,err
	}
	return publickey,nil
}

// 加密
func RSAEncrypt(plaintext string,publickey *rsa.PublicKey) (string,error) {
	cipherbytes,err := rsa.EncryptPKCS1v15(rand.Reader,publickey,[]byte(plaintext))
	ciphertext := base64.StdEncoding.EncodeToString(cipherbytes)
	return ciphertext,err
}

// 解密
func RSADecrypt(ciphertxt string,privatekey *rsa.PrivateKey) (string,error) {
	cipherbytes,err := base64.StdEncoding.DecodeString(ciphertxt)
	if err != nil {
		return "",err
	}
	publicbytes,err := rsa.DecryptPKCS1v15(rand.Reader,privatekey,cipherbytes)
	if err != nil {
		return "",err
	}
	return string(publicbytes),nil
}