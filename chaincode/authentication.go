package main

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
	"math/big"
	"time"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

type CertInformation struct {
	Country            []string
	Organization       []string
	OrganizationalUnit []string
	EmailAddress       []string
	Province           []string
	StreetAddress      []string
	SubjectKeyId       []byte
	Locality           []string
}

// Generate PrivateKey and RootCert
func (sc *SuperChain) GeneratePrivateKeyAndRootCert()(string, string){

	// cert info
	certInfo := CertInformation{
		Country:            []string{"China"},
		Organization:       []string{"buaa"},
		OrganizationalUnit: []string{"www.buaa.edu.cn"},
		EmailAddress:       []string{"wlkjaq@buaa.edu.cn"},
		StreetAddress:      []string{"37"},
		Province:           []string{"Beijing"},
		Locality:           []string{"haidian"},
		SubjectKeyId:       []byte{6, 5, 4, 3, 2, 1},
	}
	// use certInformation to new cert
	certTemp := sc.newCertificate(certInfo)
	// generate private key and public key
	rootPrivateKey, _ := ecdsa.GenerateKey(elliptic.P521(), rand.Reader)
	// create root cert []byte
	rootCertificateByte, _ := x509.CreateCertificate(rand.Reader, certTemp, certTemp, &rootPrivateKey.PublicKey, rootPrivateKey)
	// get private key []bytes
	rootPrivateKeyByte, _ := x509.MarshalECPrivateKey(rootPrivateKey)
	// get private key string and public key string
	rootCertificateString:= string(rootCertificateByte)
	rootPrivateKeyByteString := string(rootPrivateKeyByte)

	return rootCertificateString, rootPrivateKeyByteString
}

// Generate PrivateKey and RootCert
func (sc *SuperChain) CreateCertWithCsr(stub shim.ChaincodeStubInterface, csrByte []byte)([]byte, error){

	var err error

	csr, err := x509.ParseCertificateRequest(csrByte)
	if err != nil {
		return nil, err
	}
	certTemp := sc.newCertificateWithCSR(csr)
	// get root certificate
	rootCertificateByte, err := stub.GetState(RootCertificate)
	if err != nil {
		return nil, err
	}
	rootCertificate, err := x509.ParseCertificate(rootCertificateByte)
	if err != nil {
		return nil, err
	}
	// get root private key
	rootPrivateKeyByte, err := stub.GetState(RootPrivateKey)
	if err != nil {
		return nil, err
	}
	rootPrivateKey, err := x509.ParseECPrivateKey(rootPrivateKeyByte)
	if err != nil {
		return nil, err
	}
	// create new certificate
	newCertByte, err := x509.CreateCertificate(rand.Reader, certTemp, rootCertificate, csr.PublicKey, rootPrivateKey)
	if err != nil {
		return nil, err
	}

	return newCertByte, err
}

// new cert with info
func (sc *SuperChain) newCertificate(info CertInformation) *x509.Certificate {
	return &x509.Certificate{
		SerialNumber: big.NewInt(1653),
		Subject: pkix.Name{
			Country:            info.Country,
			Organization:       info.Organization,
			OrganizationalUnit: info.OrganizationalUnit,
			Province:           info.Province,
			StreetAddress:      info.StreetAddress,
		},
		NotBefore: time.Now(),
		NotAfter:  time.Now().AddDate(10, 0, 0),
		//SubjectKeyId:          info.SubjectKeyId,
		BasicConstraintsValid: true,
		IsCA:                  true,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
		//EmailAddresses: info.EmailAddress,
	}
}

// new cert with csr
func (sc *SuperChain) newCertificateWithCSR(req *x509.CertificateRequest) *x509.Certificate {
	return &x509.Certificate{
		SerialNumber: big.NewInt(1653),
		Subject:      req.Subject,
		NotBefore:    time.Now(),
		NotAfter:     time.Now().AddDate(10, 0, 0),
		//SubjectKeyId:          info.SubjectKeyId,
		BasicConstraintsValid: true,
		IsCA:                  true,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
		//EmailAddresses: info.EmailAddress,
	}
}