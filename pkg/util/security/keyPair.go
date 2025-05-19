package security

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/asn1"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"fmt"
	"log"
	"time"
)

// ParseCertificate 解析证书
func ParseCertificate(cert string) (*x509.Certificate, error) {
	block, _ := pem.Decode([]byte(cert))
	if block == nil {
		return nil, errors.New("pem.Decode error")
	}

	return x509.ParseCertificate(block.Bytes)
}

func ParseCertPublicKey(cert string) (pub *rsa.PublicKey, certID string, err error) {
	block, _ := pem.Decode([]byte(cert))
	if block == nil {
		err = errors.New("pem.Decode error")
		return
	}

	certificate, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return
	}

	pub = certificate.PublicKey.(*rsa.PublicKey)
	certID = certificate.SerialNumber.String()
	return
}

func ParsePKCS1PrivateKey(sk string) (*rsa.PrivateKey, error) {
	block, _ := pem.Decode([]byte(sk))
	if block == nil {
		return nil, errors.New("pem.Decode error")
	}
	key, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	return key, nil
}

func ParsePKCS8PrivateKey(sk string) (*rsa.PrivateKey, error) {
	block, _ := pem.Decode([]byte(sk))
	if block == nil {
		return nil, errors.New("pem.Decode error")
	}
	keyInterface, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	key := keyInterface.(*rsa.PrivateKey)

	return key, nil
}

// MarshalPKCS8PrivateKey MarshalPKCS8PrivateKey
func MarshalPKCS8PrivateKey(key *rsa.PrivateKey) []byte {
	info := struct {
		Version             int
		PrivateKeyAlgorithm []asn1.ObjectIdentifier
		PrivateKey          []byte
	}{}
	info.Version = 0
	info.PrivateKeyAlgorithm = make([]asn1.ObjectIdentifier, 1)
	info.PrivateKeyAlgorithm[0] = asn1.ObjectIdentifier{1, 2, 840, 113549, 1, 1, 1}
	info.PrivateKey = x509.MarshalPKCS1PrivateKey(key)

	k, err := asn1.Marshal(info)
	if err != nil {
		log.Panic(err.Error())
	}
	return k
}

// GenerateRSAKeyDer 生成RSA密钥对Der格式
func GenerateRSAKeyDer() (privDer, pubDer []byte, err error) {
	/* Shamelessly borrowed and adapted from some golang-samples */
	priv, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, nil, err
	}
	if err := priv.Validate(); err != nil {
		return nil, nil, fmt.Errorf("RSA key validation failed: %s", err)
	}
	privDer = MarshalPKCS8PrivateKey(priv)
	pub := priv.PublicKey
	pubDer, err = x509.MarshalPKIXPublicKey(&pub)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get der format for public key: %v", err)
	}
	return privDer, pubDer, nil
}

// GenerateRSACertDer 生成私钥、证书对
func GenerateRSACertDer(expiry time.Duration) (privDer, certDer []byte, certID string, err error) {

	priv, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return
	}
	if err = priv.Validate(); err != nil {
		return
	}
	privDer = MarshalPKCS8PrivateKey(priv)

	certDer, certID, err = CreatCertificateDer(priv, expiry)
	if err != nil {
		return
	}

	return
}

// GenerateRSACertPem 生成私钥、证书对
func GenerateRSACertPem(expiry time.Duration) (privPem, certPem, certID string, err error) {
	privDer, certDer, certID, err := GenerateRSACertDer(expiry)
	if err != nil {
		return
	}
	privPem = string(pem.EncodeToMemory(&pem.Block{Type: "PRIVATE KEY", Bytes: privDer}))
	certPem = string(pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: certDer}))
	return
}

// GenerateRSAKeysPem 生成RSA密钥对
func GenerateRSAKeysPem() (privPem, pubPem string, err error) {
	privDer, pubDer, err := GenerateRSAKeyDer()
	if err != nil {
		return "", "", err
	}

	privBlk := pem.Block{
		Type:    "RSA PRIVATE KEY",
		Headers: nil,
		Bytes:   privDer,
	}
	privPem = string(pem.EncodeToMemory(&privBlk))

	pubBlk := pem.Block{
		Type:    "PUBLIC KEY",
		Headers: nil,
		Bytes:   pubDer,
	}
	pubPem = string(pem.EncodeToMemory(&pubBlk))
	return privPem, pubPem, nil
}

// Base64KeyToPem  密钥：base64转pem格式
func Base64KeyToPem(skType, sk string) (string, error) {
	val, err := base64.StdEncoding.DecodeString(sk)
	if val == nil {
		return "", fmt.Errorf("decodeString error:%v", err)
	}

	/* For some reason chef doesn't label the keys RSA PRIVATE/PUBLIC KEY */
	privBlk := pem.Block{
		Type:    skType,
		Headers: nil,
		Bytes:   val,
	}
	privPem := string(pem.EncodeToMemory(&privBlk))

	return privPem, nil
}

// PemTobase64 密钥：pem格式转base64
func PemTobase64(sk string) (string, error) {
	block, _ := pem.Decode([]byte(sk))
	if block == nil {
		return "", errors.New("pem.Decode error")
	}

	return base64.StdEncoding.EncodeToString(block.Bytes), nil
}

// DerTobase64 密钥：der格式转base64
// base64.StdEncoding.EncodeToString(bytes)

// CertPublicKeyToPem cert->证书
func CertPublicKeyToPem(cert string) (string, error) {
	pub, _, err := ParseCertPublicKey(cert)
	if err != nil {
		return "", err
	}

	asn1Bytes, err := asn1.Marshal(*pub)
	if err != nil {
		return "", err
	}

	pubBlk := pem.Block{
		Type:    "PUBLIC KEY",
		Headers: nil,
		Bytes:   asn1Bytes,
	}
	pubPem := string(pem.EncodeToMemory(&pubBlk))
	_, err = x509.ParsePKIXPublicKey(asn1Bytes)
	if err != nil {
		return "", err
	}

	return pubPem, nil
}
