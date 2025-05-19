package security

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"errors"
	"math/big"
	"time"
)

// CreatCertificatePem 根据rsa.PrivateKey生成证书
func CreatCertificatePem(priKey *rsa.PrivateKey, expiry time.Duration) (certPem []byte, err error) {
	certBytes, _, err := CreatCertificateDer(priKey, expiry)
	if err != nil {
		return nil, err
	}

	certBlk := pem.Block{
		Type:    "CERTIFICATE",
		Headers: nil,
		Bytes:   certBytes,
	}
	return pem.EncodeToMemory(&certBlk), nil
}

// CreatCertificateDer 根据rsa.PrivateKey生成证书
func CreatCertificateDer(priKey *rsa.PrivateKey, expiry time.Duration) (certDer []byte, certID string, err error) {
	now := time.Now()
	serialNumber, _ := new(big.Int).SetString(now.Format("20060102150405"), 10)
	if serialNumber == nil {
		err = errors.New("big int format err")
		return
	}
	notBefore := now.Add(-5 * time.Minute).UTC()
	template := x509.Certificate{
		SerialNumber:          serialNumber,
		NotBefore:             notBefore,
		NotAfter:              notBefore.Add(expiry).UTC(),
		BasicConstraintsValid: true,
		Subject: pkix.Name{
			Country:      []string{"CN"},
			Province:     []string{"Shanghai"},
			Organization: []string{"Cardinfolink"},
			CommonName:   "everonet.com",
		},
	}

	certDer, err = x509.CreateCertificate(rand.Reader, &template, &template, priKey.Public(), priKey)
	if err != nil {
		return nil, "", err
	}
	certID = serialNumber.String()
	return
}
