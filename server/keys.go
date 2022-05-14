package server

import (
	"crypto/ed25519"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/mikesmitty/edkey"
	"golang.org/x/crypto/ssh"
)

func generateRSAKey() (*rsa.PrivateKey, ssh.PublicKey, error) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		return nil, nil, err
	}

	publicKey, err := ssh.NewPublicKey(&privateKey.PublicKey)
	if err != nil {
		return nil, nil, err
	}

	return privateKey, publicKey, nil
}

func generateEd25519Key() (*ed25519.PrivateKey, ssh.PublicKey, error) {
	rawPublicKey, privateKey, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		return nil, nil, err
	}

	publicKey, err := ssh.NewPublicKey(rawPublicKey)
	if err != nil {
		return nil, nil, err
	}

	return &privateKey, publicKey, nil
}

func convertRSAPrivateKeyToPem(privateKey *rsa.PrivateKey) []byte {
	return pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
	})
}

func convertEd25519PrivateKeyToPem(privateKey *ed25519.PrivateKey) []byte {
	return pem.EncodeToMemory(&pem.Block{
		Type:  "OPENSSH PRIVATE KEY",
		Bytes: edkey.MarshalED25519PrivateKey(*privateKey),
	})
}

func generateCertificate(publicKey ssh.PublicKey, subject string, certType uint32) *ssh.Certificate {
	permissions := ssh.Permissions{
		CriticalOptions: map[string]string{},
		Extensions: map[string]string{
			"permit-agent-forwarding": "",
		},
	}

	return &ssh.Certificate{
		Serial:          0,
		CertType:        certType,
		Key:             publicKey,
		KeyId:           uuid.New().String(),
		ValidPrincipals: []string{subject},
		Permissions:     permissions,
		ValidAfter:      uint64(time.Now().Unix()),
		ValidBefore:     uint64(time.Now().Unix() + 3600),
	}
}

func NewSignedRSAKeypair(caPrivateKey *rsa.PrivateKey, subject string, certType uint32) (*rsa.PrivateKey, *ssh.Certificate, error) {
	userPrivateKey, userPublicKey, err := generateRSAKey()
	if err != nil {
		return nil, nil, err
	}

	signer, err := ssh.NewSignerFromKey(caPrivateKey)
	if err != nil {
		return nil, nil, err
	}

	certificate := generateCertificate(userPublicKey, subject, certType)

	return userPrivateKey, certificate, certificate.SignCert(rand.Reader, signer)
}

func NewSignedEd25519Keypair(caPrivateKey *ed25519.PrivateKey, subject string, certType uint32) (*ed25519.PrivateKey, *ssh.Certificate, error) {
	userPrivateKey, userPublicKey, err := generateEd25519Key()
	if err != nil {
		return nil, nil, err
	}

	signer, err := ssh.NewSignerFromKey(caPrivateKey)
	if err != nil {
		return nil, nil, err
	}

	certificate := generateCertificate(userPublicKey, subject, certType)

	return userPrivateKey, certificate, certificate.SignCert(rand.Reader, signer)
}
