package server

import (
	"crypto/ed25519"
	"crypto/rsa"
	"time"

	"golang.org/x/crypto/ssh"
	"gorm.io/gorm"

	"github.com/mfinelli/rush/db"
)

func getLatestCA(t string, rdb *gorm.DB) (interface{}, string, *db.CACertificate, error) {
	var ca db.CACertificate
	err := rdb.Where("type = ?", t).Order("updated_at desc").First(&ca).Error
	if err != nil {
		return nil, "", &db.CACertificate{}, err
	}

	caPrivateKey, err := ssh.ParseRawPrivateKey([]byte(ca.PrivateKey))
	if err != nil {
		return nil, "", &db.CACertificate{}, err
	}

	return caPrivateKey, ca.PublicKey, &ca, nil
}

func getLatestEd25519CA(rdb *gorm.DB) (*ed25519.PrivateKey, string, *db.CACertificate, error) {
	caPrivateKey, caPublicKey, ca, err := getLatestCA("ed25519", rdb)
	if err != nil {
		return nil, "", &db.CACertificate{}, err
	}

	return caPrivateKey.(*ed25519.PrivateKey), caPublicKey, ca, nil
}

func getLatestRSACA(rdb *gorm.DB) (*rsa.PrivateKey, string, *db.CACertificate, error) {
	caPrivateKey, caPublicKey, ca, err := getLatestCA("rsa", rdb)
	if err != nil {
		return nil, "", &db.CACertificate{}, err
	}

	return caPrivateKey.(*rsa.PrivateKey), caPublicKey, ca, nil
}

func saveHostKey(rdb *gorm.DB, t string, cert *ssh.Certificate, signedKey string, privateKey string, signer *db.CACertificate) error {
	// TODO: error handling
	rdb.Create(&db.HostCertificate{
		KeyId:           cert.KeyId,
		Hostname:        cert.ValidPrincipals[0],
		Type:            t,
		NotBefore:       time.Unix(int64(cert.ValidAfter), 0),
		NotAfter:        time.Unix(int64(cert.ValidBefore), 0),
		PublicKey:       string(ssh.MarshalAuthorizedKey(cert.Key)),
		SignedPublicKey: signedKey,
		PrivateKey:      privateKey,
		CACertificate:   *signer,
	})

	return nil
}

func saveUserKey(rdb *gorm.DB, t string, cert *ssh.Certificate, signedKey string, privateKey string, signer *db.CACertificate) error {
	// TODO: error handling
	rdb.Create(&db.UserCertificate{
		KeyId:           cert.KeyId,
		Username:        cert.ValidPrincipals[0],
		Type:            t,
		NotBefore:       time.Unix(int64(cert.ValidAfter), 0),
		NotAfter:        time.Unix(int64(cert.ValidBefore), 0),
		PublicKey:       string(ssh.MarshalAuthorizedKey(cert.Key)),
		SignedPublicKey: signedKey,
		PrivateKey:      privateKey,
		CACertificate:   *signer,
	})

	return nil
}
