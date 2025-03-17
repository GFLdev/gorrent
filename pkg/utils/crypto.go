package utils

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/x509"
	"encoding/pem"
	"fmt"
)

const PrivateKeyPath = "key.rsa"
const PublicKeyPath = "key.rsa.pub"

type Keys struct {
	PrivateKey []byte
	PublicKey  []byte
}

func GenerateKeys() (Keys, error) {
	bitSize := 4096

	// Generate RSA key
	key, err := rsa.GenerateKey(rand.Reader, bitSize)
	if err != nil {
		return Keys{}, fmt.Errorf("could not generate key: %w", err)
	}

	// Encode private key and public key to PKCS#1 ASN.1 PEM
	keyPEM := pem.EncodeToMemory(
		&pem.Block{
			Type:  "RSA PRIVATE KEY",
			Bytes: x509.MarshalPKCS1PrivateKey(key),
		},
	)
	pubPEM := pem.EncodeToMemory(
		&pem.Block{
			Type:  "PUBLIC KEY",
			Bytes: x509.MarshalPKCS1PublicKey(&key.PublicKey),
		},
	)

	// Save to disk
	err = WriteFile(PrivateKeyPath, keyPEM)
	if err != nil {
		return Keys{}, fmt.Errorf("could not save private key to disk: %w", err)
	}
	err = WriteFile(PublicKeyPath, pubPEM)
	if err != nil {
		return Keys{}, fmt.Errorf("could not save public key to disk: %w", err)
	}
	return Keys{PrivateKey: keyPEM, PublicKey: pubPEM}, nil
}

func SHA1Encode(val []byte) []byte {
	h := sha1.New()
	h.Write(val)
	return h.Sum(nil)
}
