package main

import (
	"crypto/ed25519"
	"crypto/rand"
	"encoding/pem"
	"fmt"
	"golang.org/x/crypto/ssh"
	"log"
	"time"
)

import "github.com/mikesmitty/edkey"
import "github.com/google/uuid"

func main() {
	permissions := ssh.Permissions{
		CriticalOptions: map[string]string{},
		Extensions:      map[string]string{"permit-agent-forwarding": ""},
	}

	rawPublicKey, rawPrivateKey, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		log.Fatal(err.Error())
	}

	sshPublicKey, err := ssh.NewPublicKey(rawPublicKey)
	if err != nil {
		log.Fatal(err.Error())
	}

	pemPrivateKey := &pem.Block{
		Type:  "OPENSSH PRIVATE KEY",
		Bytes: edkey.MarshalED25519PrivateKey(rawPrivateKey),
	}

	privateKey := pem.EncodeToMemory(pemPrivateKey)
	publicKey := ssh.MarshalAuthorizedKey(sshPublicKey)

	fmt.Println(string(privateKey))
	fmt.Println(string(publicKey))

	rawUserPublicKey, _, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		log.Fatal(err.Error())
	}

	sshUserPublicKey, err := ssh.NewPublicKey(rawUserPublicKey)
	if err != nil {
		log.Fatal(err.Error())
	}

	cert := ssh.Certificate{
		Serial:          0,
		CertType:        ssh.UserCert,
		Key:             sshUserPublicKey,
		KeyId:           uuid.New().String(),
		ValidPrincipals: []string{"test"},
		Permissions:     permissions,
		ValidAfter:      uint64(time.Now().Unix()),
		ValidBefore:     uint64(time.Now().Unix() + 30),
	}

	// fmt.Printf("%v\n", cert)

	signer, err := ssh.NewSignerFromKey(rawPrivateKey)
	if err != nil {
		log.Fatal(err.Error())
	}

	err = cert.SignCert(rand.Reader, signer)
	if err != nil {
		log.Fatal(err.Error())
	}

	// fmt.Printf("%v\n", cert)
	signedKey := ssh.MarshalAuthorizedKey(&cert)
	fmt.Println(string(ssh.MarshalAuthorizedKey(cert.SignatureKey)))
	fmt.Println(string(signedKey))

}
