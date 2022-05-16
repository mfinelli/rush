package server

import (
	"crypto/ed25519"
	"crypto/rsa"

	"golang.org/x/crypto/ssh"
	"gorm.io/gorm"

	"github.com/mfinelli/rush/db"
)

type CACertificateResponse struct {
	PublicKey string `json:"public_key"`
}

type SignedKeyResponse struct {
	CA         string `json:"ca"`
	PublicKey  string `json:"public_key"`
	PrivateKey string `json:"private_key"`
}

func getLatestCA(t string, rdb *gorm.DB) (interface{}, string, error) {
	var ca db.CACertificate
	err := rdb.Where("type = ?", t).Order("updated_at desc").First(&ca).Error
	if err != nil {
		return nil, "", err
	}

	caPrivateKey, err := ssh.ParseRawPrivateKey([]byte(ca.PrivateKey))
	if err != nil {
		return nil, "", err
	}

	return caPrivateKey, ca.PublicKey, nil
}

func getLatestEd25519CA(rdb *gorm.DB) (*ed25519.PrivateKey, string, error) {
	caPrivateKey, caPublicKey, err := getLatestCA("ed25519", rdb)
	if err != nil {
		return nil, "", err
	}

	return caPrivateKey.(*ed25519.PrivateKey), caPublicKey, nil
}

func getLatestRSACA(rdb *gorm.DB) (*rsa.PrivateKey, string, error) {
	caPrivateKey, caPublicKey, err := getLatestCA("rsa", rdb)
	if err != nil {
		return nil, "", err
	}

	return caPrivateKey.(*rsa.PrivateKey), caPublicKey, nil
}

func generateRSAUserKey(rdb *gorm.DB, username string) (SignedKeyResponse, error) {
	caPrivateKey, caPublicKey, err := getLatestRSACA(rdb)
	if err != nil {
		return SignedKeyResponse{}, err
	}

	privateKey, signedCertificate, err := NewSignedRSAKeypair(caPrivateKey, username, ssh.UserCert)
	if err != nil {
		return SignedKeyResponse{}, err
	}

	return SignedKeyResponse{
		CA:         caPublicKey,
		PublicKey:  string(ssh.MarshalAuthorizedKey(signedCertificate)),
		PrivateKey: string(convertRSAPrivateKeyToPem(privateKey)),
	}, nil
}

func generateEd25519UserKey(rdb *gorm.DB, username string) (SignedKeyResponse, error) {
	caPrivateKey, caPublicKey, err := getLatestEd25519CA(rdb)
	if err != nil {
		return SignedKeyResponse{}, err
	}

	privateKey, signedCertificate, err := NewSignedEd25519Keypair(caPrivateKey, username, ssh.UserCert)
	if err != nil {
		return SignedKeyResponse{}, err
	}

	return SignedKeyResponse{
		CA:         caPublicKey,
		PublicKey:  string(ssh.MarshalAuthorizedKey(signedCertificate)),
		PrivateKey: string(convertEd25519PrivateKeyToPem(privateKey)),
	}, nil
}

func generateRSAHostKey(rdb *gorm.DB, hostname string) (SignedKeyResponse, error) {
	caPrivateKey, caPublicKey, err := getLatestRSACA(rdb)
	if err != nil {
		return SignedKeyResponse{}, err
	}

	privateKey, signedCertificate, err := NewSignedRSAKeypair(caPrivateKey, hostname, ssh.HostCert)
	if err != nil {
		return SignedKeyResponse{}, err
	}

	return SignedKeyResponse{
		CA:         caPublicKey,
		PublicKey:  string(ssh.MarshalAuthorizedKey(signedCertificate)),
		PrivateKey: string(convertRSAPrivateKeyToPem(privateKey)),
	}, nil
}

func generateEd25519HostKey(rdb *gorm.DB, hostname string) (SignedKeyResponse, error) {
	caPrivateKey, caPublicKey, err := getLatestEd25519CA(rdb)
	if err != nil {
		return SignedKeyResponse{}, err
	}

	privateKey, signedCertificate, err := NewSignedEd25519Keypair(caPrivateKey, hostname, ssh.HostCert)
	if err != nil {
		return SignedKeyResponse{}, err
	}

	return SignedKeyResponse{
		CA:         caPublicKey,
		PublicKey:  string(ssh.MarshalAuthorizedKey(signedCertificate)),
		PrivateKey: string(convertEd25519PrivateKeyToPem(privateKey)),
	}, nil
}
