package server

import (
	"golang.org/x/crypto/ssh"
	"gorm.io/gorm"
)

type CACertificateResponse struct {
	PublicKey string `json:"public_key"`
}

type SignedKeyResponse struct {
	CA         string `json:"ca"`
	PublicKey  string `json:"public_key"`
	PrivateKey string `json:"private_key"`
}

func generateRSAUserKey(rdb *gorm.DB, username string) (SignedKeyResponse, error) {
	caPrivateKey, caPublicKey, ca, err := getLatestRSACA(rdb)
	if err != nil {
		return SignedKeyResponse{}, err
	}

	privateKey, signedCertificate, err := NewSignedRSAKeypair(caPrivateKey, username, ssh.UserCert)
	if err != nil {
		return SignedKeyResponse{}, err
	}

	pemPublicKey := string(ssh.MarshalAuthorizedKey(signedCertificate))
	pemPrivateKey := string(convertRSAPrivateKeyToPem(privateKey))

	err = saveUserKey(rdb, "rsa", signedCertificate, pemPublicKey, pemPrivateKey, ca)
	if err != nil {
		return SignedKeyResponse{}, err
	}

	return SignedKeyResponse{
		CA:         caPublicKey,
		PublicKey:  pemPublicKey,
		PrivateKey: pemPrivateKey,
	}, nil
}

func generateEd25519UserKey(rdb *gorm.DB, username string) (SignedKeyResponse, error) {
	caPrivateKey, caPublicKey, ca, err := getLatestEd25519CA(rdb)
	if err != nil {
		return SignedKeyResponse{}, err
	}

	privateKey, signedCertificate, err := NewSignedEd25519Keypair(caPrivateKey, username, ssh.UserCert)
	if err != nil {
		return SignedKeyResponse{}, err
	}

	pemPublicKey := string(ssh.MarshalAuthorizedKey(signedCertificate))
	pemPrivateKey := string(convertEd25519PrivateKeyToPem(privateKey))

	err = saveUserKey(rdb, "ed25519", signedCertificate, pemPublicKey, pemPrivateKey, ca)
	if err != nil {
		return SignedKeyResponse{}, err
	}

	return SignedKeyResponse{
		CA:         caPublicKey,
		PublicKey:  pemPublicKey,
		PrivateKey: pemPrivateKey,
	}, nil
}

func generateRSAHostKey(rdb *gorm.DB, hostname string) (SignedKeyResponse, error) {
	caPrivateKey, caPublicKey, ca, err := getLatestRSACA(rdb)
	if err != nil {
		return SignedKeyResponse{}, err
	}

	privateKey, signedCertificate, err := NewSignedRSAKeypair(caPrivateKey, hostname, ssh.HostCert)
	if err != nil {
		return SignedKeyResponse{}, err
	}

	pemPublicKey := string(ssh.MarshalAuthorizedKey(signedCertificate))
	pemPrivateKey := string(convertRSAPrivateKeyToPem(privateKey))

	err = saveHostKey(rdb, "rsa", signedCertificate, pemPublicKey, pemPrivateKey, ca)
	if err != nil {
		return SignedKeyResponse{}, err
	}

	return SignedKeyResponse{
		CA:         caPublicKey,
		PublicKey:  pemPublicKey,
		PrivateKey: pemPrivateKey,
	}, nil
}

func generateEd25519HostKey(rdb *gorm.DB, hostname string) (SignedKeyResponse, error) {
	caPrivateKey, caPublicKey, ca, err := getLatestEd25519CA(rdb)
	if err != nil {
		return SignedKeyResponse{}, err
	}

	privateKey, signedCertificate, err := NewSignedEd25519Keypair(caPrivateKey, hostname, ssh.HostCert)
	if err != nil {
		return SignedKeyResponse{}, err
	}

	pemPublicKey := string(ssh.MarshalAuthorizedKey(signedCertificate))
	pemPrivateKey := string(convertEd25519PrivateKeyToPem(privateKey))

	err = saveHostKey(rdb, "ed25519", signedCertificate, pemPublicKey, pemPrivateKey, ca)
	if err != nil {
		return SignedKeyResponse{}, err
	}

	return SignedKeyResponse{
		CA:         caPublicKey,
		PublicKey:  pemPublicKey,
		PrivateKey: pemPrivateKey,
	}, nil
}
