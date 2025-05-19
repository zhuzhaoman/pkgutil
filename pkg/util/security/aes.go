package security

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"github.com/zhuzhaoman/pkgutil/pkg/log"
	"io"
	"strings"
)

// AESCBCMode 如果key位base64编码过的字符串
type AESCBCMode struct {
	Key    []byte
	Err    error
	sysAES *AESCBCMode
}

// NewAESCBCEncrypt 创建一个 AES 加密对象，使用 CBC 模式
func NewAESCBCEncrypt(b64Key, sysKey string) *AESCBCMode {
	bytesKey, err := base64.StdEncoding.DecodeString(b64Key)

	if err != nil {
		log.Errorf("AES key(%s) base64 decode error: %s", b64Key, err)
	}

	return &AESCBCMode{
		Key: bytesKey,
		Err: err,
		sysAES: &AESCBCMode{
			Key: []byte(sysKey),
		},
	}
}

// DcyAndUseSysKeyEcy
// decrypted 解密后的明文 encrypted 使用新key后的密文
func (a *AESCBCMode) DcyAndUseSysKeyEcy(ct, sysKey string) (decrypted, encrypted string) {
	if a.sysAES == nil {
		a.sysAES = &AESCBCMode{Key: []byte(sysKey)}
	}

	// decrypt
	decrypted = a.Decrypt(ct)

	if a.Err != nil {
		log.Errorf("Decrypt the string error, string is %s, error is %s", ct, a.Err.Error())
		return decrypted, decrypted
	}
	// encrypt
	encrypted = a.sysAES.Encrypt(decrypted)

	if a.sysAES.Err != nil {
		// 将错误传递到a
		a.Err = a.sysAES.Err
		log.Errorf("Encrypt the string error, string is %s, error is %s", decrypted, a.Err.Error())
	}
	return
}

func (a *AESCBCMode) AesCbcEncrypt(ct, sysKey string) string {
	if a.sysAES == nil {
		a.sysAES = &AESCBCMode{Key: []byte(sysKey)}
	}

	encrypted := a.sysAES.Encrypt(ct)

	if a.sysAES.Err != nil {
		return a.sysAES.Err.Error()
	}
	return encrypted
}

// Encrypt cbc mode
func (a *AESCBCMode) Encrypt(pt string) string {
	if a.Err != nil {
		return pt
	}
	plaintext := pkcs7Padding([]byte(pt), aes.BlockSize)

	if len(plaintext)%aes.BlockSize != 0 {
		a.Err = fmt.Errorf("%s : plaintext is not a multiple of the block size", pt)
		log.Error(a.Err.Error())
		return pt
	}

	block, err := aes.NewCipher(a.Key)
	if err != nil {
		a.Err = err
		log.Error(a.Err.Error())
		return pt
	}

	ciphertext := make([]byte, aes.BlockSize+len(plaintext))
	iv := ciphertext[:aes.BlockSize]
	// 随机生成16个字节数组
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		a.Err = err
		log.Errorf("generate the rand num error, error is %s", a.Err.Error())
		return pt
	}
	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(ciphertext[aes.BlockSize:], plaintext)

	return base64.StdEncoding.EncodeToString(ciphertext)
}

// Decrypt cbc mode
func (a *AESCBCMode) Decrypt(ct string) string {

	if a.Err != nil {
		return ct
	}
	defer func() {
		if err := recover(); err != nil {
			a.Err = err.(error)
		}
	}()

	ct = strings.TrimSpace(ct)
	ciphertext, err := base64.StdEncoding.DecodeString(ct)
	if err != nil {
		a.Err = err
		log.Errorf("Decode the string error, the string is %s, error is %s", ct, err.Error())
		return ct
	}
	block, err := aes.NewCipher(a.Key)
	if err != nil {
		a.Err = err
		log.Errorf("create the block error, the error is %s", err.Error())
		return ct
	}

	if len(ciphertext) < aes.BlockSize {
		a.Err = fmt.Errorf("%s : ciphertext too short", ct)
		log.Error(a.Err.Error())
		return ct
	}
	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]

	if len(ciphertext)%aes.BlockSize != 0 {
		a.Err = fmt.Errorf("%s : ciphertext is not a multiple of the block size", ct)
		log.Error(a.Err.Error())
		return ct
	}

	mode := cipher.NewCBCDecrypter(block, iv)

	mode.CryptBlocks(ciphertext, ciphertext)
	ciphertext = pkcs7UnPadding(ciphertext)
	if ciphertext == nil {
		a.Err = fmt.Errorf("the ciphertext is nil")
		log.Error(a.Err.Error())
		return ct
	}
	return string(ciphertext)
}

func pkcs7Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

func pkcs7UnPadding(origData []byte) []byte {
	length := len(origData)
	if length == 0 {
		log.Error("the length of the origData is empty")
		return nil
	}

	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}

// AESCFBMode CFB 模式的 AES 加密
type AESCFBMode struct {
	Key []byte
	Err error
}

// NewAESCFBEncrypt 创建AES加密处理对象
func NewAESCFBEncrypt(sKey string) *AESCFBMode {

	return &AESCFBMode{
		Key: []byte(sKey),
		Err: nil,
	}
}

// Encrypt aesCFBEncrypt aes 加密  对商户敏感信息加密
func (a *AESCFBMode) Encrypt(pt string) string {

	block, err := aes.NewCipher(a.Key)
	if err != nil {
		a.Err = err
		log.Error("create the block error")
		return ""
	}
	plaintext := []byte(pt)
	ciphertext := make([]byte, aes.BlockSize+len(plaintext))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		a.Err = err
		log.Error("generate the rand num error")
		return ""
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	//stream := cipher.NewCFBEncrypter(block, []byte("0000000000000000"))
	stream.XORKeyStream(ciphertext[aes.BlockSize:], plaintext)

	return hex.EncodeToString(ciphertext)
}

// Decrypt aes 解密
func (a *AESCFBMode) Decrypt(ct string) string {

	block, err := aes.NewCipher(a.Key)
	if err != nil {
		a.Err = err
		log.Error("create the block error")
		return ""
	}

	ciphertext, err := hex.DecodeString(ct)
	if err != nil {
		a.Err = err
		log.Errorf("Decode string error, the string is %s, error is %s", ct, err.Error())
		return ""
	}

	if len(ciphertext) < aes.BlockSize {
		err = fmt.Errorf("the length of the ciphertext is too short")
		a.Err = err
		log.Errorf("%s", err)
		return ""
	}
	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(ciphertext, ciphertext)
	return string(ciphertext)
}

// AESEncryptECBStr ECB加密
func AESEncryptECBStr(source string, keys string) string {
	if source == "" {
		return ""
	}
	// 字符串转换成切片
	src := []byte(source)
	key := []byte(keys)
	cipher, _ := aes.NewCipher(generateKeys(key))
	length := (len(src) + aes.BlockSize) / aes.BlockSize
	plain := make([]byte, length*aes.BlockSize)
	copy(plain, src)
	pad := byte(len(plain) - len(src))
	for i := len(src); i < len(plain); i++ {
		plain[i] = pad
	}
	encrypted := make([]byte, len(plain))
	// 分组分块加密
	for bs, be := 0, cipher.BlockSize(); bs <= len(src); bs, be = bs+cipher.BlockSize(), be+cipher.BlockSize() {
		cipher.Encrypt(encrypted[bs:be], plain[bs:be])
	}

	//return encrypted
	encryptstr := strings.ToUpper(hex.EncodeToString(encrypted))
	return encryptstr
}

// AESDecryptECBStr ECB解密
func AESDecryptECBStr(encrypteds string, keys string) string {
	if encrypteds == "" {
		return ""
	}
	// 字符串转换成切片
	//encrypted := []byte(encrypteds)
	encrypted, _ := hex.DecodeString(encrypteds)
	key := []byte(keys)

	cipher, _ := aes.NewCipher(generateKeys(key))
	decrypted := make([]byte, len(encrypted))
	//
	for bs, be := 0, cipher.BlockSize(); bs < len(encrypted); bs, be = bs+cipher.BlockSize(), be+cipher.BlockSize() {
		cipher.Decrypt(decrypted[bs:be], encrypted[bs:be])
	}

	trim := 0
	if len(decrypted) > 0 {
		trim = len(decrypted) - int(decrypted[len(decrypted)-1])
	}
	return string(decrypted[:trim])
}

func generateKeys(key []byte) (genKey []byte) {
	genKey = make([]byte, 16)
	copy(genKey, key)
	for i := 16; i < len(key); {
		for j := 0; j < 16 && i < len(key); j, i = j+1, i+1 {
			genKey[j] ^= key[i]
		}
	}
	return genKey
}
