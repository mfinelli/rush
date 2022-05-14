package main

import (
	"crypto/ed25519"
	"crypto/rand"
	// "crypto/x509"
	"encoding/pem"
	"fmt"
	"golang.org/x/crypto/ssh"
	"log"
	"time"
)

import "github.com/mikesmitty/edkey"
import "github.com/google/uuid"

func main() {
	rawPublicKey, rawPrivateKey, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		log.Fatal(err.Error())
	}

	sshPublicKey, err := ssh.NewPublicKey(rawPublicKey)
	if err != nil {
		log.Fatal(err.Error())
	}

	pemPrivateKey := &pem.Block{
		Type: "OPENSSH PRIVATE KEY",
		Bytes: edkey.MarshalED25519PrivateKey(rawPrivateKey),
	}

	privateKey := pem.EncodeToMemory(pemPrivateKey)
	publicKey := ssh.MarshalAuthorizedKey(sshPublicKey)

	fmt.Println(string(privateKey))
	fmt.Println(string(publicKey))

	// https://github.com/openssh/openssh-portable/blob/master/PROTOCOL.certkeys#L73-L74
	var userType uint32 = 1
	// var hostType uint32 = 2

	cert := ssh.Certificate{
		Serial: 0,
		CertType: userType,
		Key: sshPublicKey,
		KeyId: uuid.New().String(),
		ValidPrincipals: []string{"test"},
		ValidAfter: uint64(time.Now().Unix()),
		ValidBefore: uint64(time.Now().Unix() + 30),
	}

	fmt.Printf("%v\n", cert)

	signer, err := ssh.NewSignerFromKey(rawPrivateKey)
	if err != nil {
		log.Fatal(err.Error())
	}

	err = cert.SignCert(rand.Reader, signer)
	if err != nil {
		log.Fatal(err.Error())
	}

	fmt.Printf("%v\n", cert)
}
