package server

import (
	"golang.org/x/crypto/ssh"
)

type SignedKeyResponse struct {
	CA         string `json:"ca"`
	PublicKey  string `json:"public_key"`
	PrivateKey string `json:"private_key"`
}

func generateRSAUserKey(username string) (SignedKeyResponse, error) {
	// TODO: we should be consuming a CA key, not generating it on the fly
	caPrivateKey, caPublicKey, err := generateRSAKey()
	if err != nil {
		return SignedKeyResponse{}, err
	}

	privateKey, signedCertificate, err := NewSignedRSAKeypair(caPrivateKey, username, ssh.UserCert)
	if err != nil {
		return SignedKeyResponse{}, err
	}

	return SignedKeyResponse{
		CA:         string(ssh.MarshalAuthorizedKey(caPublicKey)),
		PublicKey:  string(ssh.MarshalAuthorizedKey(signedCertificate)),
		PrivateKey: string(convertRSAPrivateKeyToPem(privateKey)),
	}, nil
}

func generateEd25519UserKey(username string) (SignedKeyResponse, error) {
	// TODO: we should be consuming a CA key, not generating it on the fly
	caPrivateKey, caPublicKey, err := generateEd25519Key()
	if err != nil {
		return SignedKeyResponse{}, err
	}

	privateKey, signedCertificate, err := NewSignedEd25519Keypair(caPrivateKey, username, ssh.UserCert)
	if err != nil {
		return SignedKeyResponse{}, err
	}

	return SignedKeyResponse{
		CA:         string(ssh.MarshalAuthorizedKey(caPublicKey)),
		PublicKey:  string(ssh.MarshalAuthorizedKey(signedCertificate)),
		PrivateKey: string(convertEd25519PrivateKeyToPem(privateKey)),
	}, nil
}

func generateRSAHostKey(hostname string) (SignedKeyResponse, error) {
	// TODO: we should be consuming a CA key, not generating it on the fly
	caPrivateKey, caPublicKey, err := generateRSAKey()
	if err != nil {
		return SignedKeyResponse{}, err
	}

	privateKey, signedCertificate, err := NewSignedRSAKeypair(caPrivateKey, hostname, ssh.HostCert)
	if err != nil {
		return SignedKeyResponse{}, err
	}

	return SignedKeyResponse{
		CA:         string(ssh.MarshalAuthorizedKey(caPublicKey)),
		PublicKey:  string(ssh.MarshalAuthorizedKey(signedCertificate)),
		PrivateKey: string(convertRSAPrivateKeyToPem(privateKey)),
	}, nil
}

func generateEd25519HostKey(hostname string) (SignedKeyResponse, error) {
	// TODO: we should be consuming a CA key, not generating it on the fly
	caPrivateKey, caPublicKey, err := generateEd25519Key()
	if err != nil {
		return SignedKeyResponse{}, err
	}

	privateKey, signedCertificate, err := NewSignedEd25519Keypair(caPrivateKey, hostname, ssh.HostCert)
	if err != nil {
		return SignedKeyResponse{}, err
	}

	return SignedKeyResponse{
		CA:         string(ssh.MarshalAuthorizedKey(caPublicKey)),
		PublicKey:  string(ssh.MarshalAuthorizedKey(signedCertificate)),
		PrivateKey: string(convertEd25519PrivateKeyToPem(privateKey)),
	}, nil
}
