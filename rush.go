package main

import (
	"crypto/ed25519"
	"crypto/rand"
	// "crypto/x509"
	"encoding/pem"
	"fmt"
	"golang.org/x/crypto/ssh"
	"log"
)

import "github.com/mikesmitty/edkey"

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
}
