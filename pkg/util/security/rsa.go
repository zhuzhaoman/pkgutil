package security

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"github.com/zhuzhaoman/pkgutil/pkg/log"
)

/*
RSA 公钥私钥 packs1有两种格式 一种是pem也就是带-----BEGIN PUBLIC KEY-----头尾这种
这种在前端jsencrypt是直接用来加密，应该是加密工具做了转换，但是在golang后端需要先去掉转义符
转换成pem后再加密
*/

func ParsePKCS1Or8PrivKey(key []byte) (*rsa.PrivateKey, error) {
	block, _ := pem.Decode(key)
	if block == nil {
		return nil, errors.New("priv key error")
	}
	var privKey *rsa.PrivateKey
	privKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err == nil {
		return privKey, nil
	}

	privKeyInterface, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err == nil {
		if priKey, ok := privKeyInterface.(*rsa.PrivateKey); ok {
			return priKey, nil
		}
	}
	return nil, errors.New("parse private key error")
}

// ParsePublicKey 解析公钥
func ParsePublicKey(key []byte) (*rsa.PublicKey, error) {
	block, _ := pem.Decode(key)
	if block == nil {
		return nil, errors.New("public key error")
	}
	publicKeyInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	if publicKey, ok := publicKeyInterface.(*rsa.PublicKey); ok {
		return publicKey, nil
	} else {
		err = errors.New("public key error")
		return nil, err
	}
}

// RSAEncrypt RSA 加密
func RSAEncrypt(origData, publicKey []byte) ([]byte, error) {
	block, _ := pem.Decode(publicKey)
	if block == nil {
		return nil, errors.New("public key error")
	}
	pubInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		log.Errorf("x509.ParsePKIXPublicKey error: %s", err)
		return nil, err
	}
	pub := pubInterface.(*rsa.PublicKey)

	return rsa.EncryptPKCS1v15(rand.Reader, pub, origData)
}

// RSAEncryptBase64 RSA 加密 Base64 密文
func RSAEncryptBase64(origData, publicKey []byte) (string, error) {
	cipherText, err := RSAEncrypt(origData, publicKey)
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(cipherText), nil
}

// RSAEncryptDERType RSA加密，密钥是 DER 格式
func RSAEncryptDERType(origData, publicDERKey []byte) ([]byte, error) {
	pubInterface, err := x509.ParsePKIXPublicKey(publicDERKey)
	if err != nil {
		log.Errorf("x509.ParsePKIXPublicKey error: %s", err)
		return nil, err
	}
	pub := pubInterface.(*rsa.PublicKey)

	return rsa.EncryptPKCS1v15(rand.Reader, pub, origData)
}

// RSADecrypt RSA 解密
func RSADecrypt(ciphertext, privateKey []byte) ([]byte, error) {
	priv, err := ParsePKCS1Or8PrivKey(privateKey)
	if err != nil {
		return nil, err
	}

	return rsa.DecryptPKCS1v15(rand.Reader, priv, ciphertext)
}

// RSADecryptBase64 RSA 解密 Base64 密文
func RSADecryptBase64(b64Cipher string, privateKey []byte) ([]byte, error) {
	cipherText, err := base64.StdEncoding.DecodeString(b64Cipher)
	if err != nil {
		log.Errorf("base64.StdEncoding.DecodeString error! b64Cipher:%s, err:%s", b64Cipher, err)
		return nil, err
	}

	return RSADecrypt(cipherText, privateKey)
}

// RSASign RSA 签名
func RSASign(text []byte, privateKey string, hashType crypto.Hash) ([]byte, error) {
	priv, err := ParsePKCS1Or8PrivKey([]byte(privateKey))
	if err != nil {
		log.Errorf("x509.ParsePKCS1PrivateKey error: %s", err)
		return nil, err
	}

	return rsa.SignPKCS1v15(rand.Reader, priv, hashType, text)
}

// RSASign RSA 签名 Base64
func RSASignBase64(text []byte, privateKey string, hashType crypto.Hash) (string, error) {
	sign, err := RSASign(text, privateKey, hashType)
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(sign), nil
}

// RSAVerifySign RSA验证签名
func RSAVerifySign(text []byte, publicKey string, hashType crypto.Hash, sig []byte) error {
	block, _ := pem.Decode([]byte(publicKey))
	if block == nil {
		log.Errorf("pem.Decode error")
		return errors.New("public key error")
	}
	pubInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		log.Errorf("x509.ParsePKIXPublicKey error: %s", err)
		return err
	}
	pub := pubInterface.(*rsa.PublicKey)
	return rsa.VerifyPKCS1v15(pub, hashType, text, sig)
}

// RSAVerifySignBase64 RSA验证Base64签名
func RSAVerifySignBase64(b64Cipher, publicKey string, hashType crypto.Hash, text []byte) error {
	cipherText, err := base64.StdEncoding.DecodeString(b64Cipher)
	if err != nil {
		cipherText, err = base64.RawURLEncoding.DecodeString(b64Cipher)
		if err != nil {
			log.Errorf("base64.StdEncoding.DecodeString error! b64Cipher:%s, err:%s", b64Cipher, err)
			return err
		}
	}

	return RSAVerifySign(text, publicKey, hashType, cipherText)
}

// SHA256WithRSABase64 SHA256WithRSA签名算法签名，返回basd64编码后的签名
func SHA256WithRSABase64(data, privateKey []byte) (string, error) {
	sign, err := SHA256WithRSA(data, privateKey, false)
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(sign), nil
}

// SHA256WithRSA SHA256WithRSA签名算法签名
func SHA256WithRSA(data, privateKey []byte, noneWithRsa bool) ([]byte, error) {
	priKey, err := ParsePKCS1Or8PrivKey(privateKey)
	if err != nil {
		return nil, err
	}

	hashed := sha256.Sum256(data)
	if noneWithRsa {
		return rsa.SignPKCS1v15(rand.Reader, priKey, crypto.Hash(0), hashed[:])
	}
	return rsa.SignPKCS1v15(rand.Reader, priKey, crypto.SHA256, hashed[:])
}

// VerifySHA256WithRSABase64 SHA256WithRSA签名算法验签，如果验签通过，则err 值为 nil
func VerifySHA256WithRSABase64(origin []byte, b64Sign, publicKey string) error {
	pubKey, err := ParsePublicKey([]byte(publicKey))
	if err != nil {
		return err
	}

	sign, err := base64.StdEncoding.DecodeString(b64Sign)
	if err != nil {
		return err
	}

	hashed := sha256.Sum256(origin)
	err = rsa.VerifyPKCS1v15(pubKey, crypto.SHA256, hashed[:], sign)

	return err
}

// key是否具有头尾换行不交由程序判断
// keyType=1,为key增加头尾，并每间隔64位换行
// ifPublic true 为公钥， false为私钥
// keyType=0,不变
func FormateKey(key string, keyType int, ifPublic bool) string {
	if keyType == 0 {
		return key
	}
	if ifPublic {
		var publicHeader = "\n-----BEGIN PUBLIC KEY-----\n"
		var publicTail = "-----END PUBLIC KEY-----\n"
		var temp string
		Split(key, &temp)
		return publicHeader + temp + publicTail
	} else {
		var publicHeader = "\n-----BEGIN RSA PRIVATE KEY-----\n"
		var publicTail = "-----END RSA PRIVATE KEY-----\n"
		var temp string
		Split(key, &temp)
		return publicHeader + temp + publicTail
	}
}

func Split(key string, temp *string) {
	if len(key) <= 64 {
		*temp = *temp + key + "\n"
	}
	for i := 0; i < len(key); i++ {
		if (i+1)%64 == 0 {
			*temp = *temp + key[:i+1] + "\n"
			log.Info(len(*temp) - 1)
			key = key[i+1:]
			Split(key, temp)
			break
		}
	}
}
