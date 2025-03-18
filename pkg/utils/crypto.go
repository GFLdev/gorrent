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

// GenerateKeys generates a new RSA key pair, encodes them in PEM format.
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
	return Keys{PrivateKey: keyPEM, PublicKey: pubPEM}, nil
}

// SaveKeys saves the provided private and public keys to the respective file paths.
func SaveKeys(keys Keys) error {
	err := WriteFile(PrivateKeyPath, keys.PrivateKey)
	if err != nil {
		return fmt.Errorf("could not save private key: %w", err)
	}
	err = WriteFile(PublicKeyPath, keys.PublicKey)
	if err != nil {
		return fmt.Errorf("could not save public key: %w", err)
	}
	return nil
}

// LoadKeys loads the public and private keys from their respective file paths.
func LoadKeys() (Keys, error) {
	privateKey, err := ReadFile(PrivateKeyPath)
	if err != nil {
		return Keys{}, fmt.Errorf("could not load private key: %w", err)
	}
	publicKey, err := ReadFile(PublicKeyPath)
	if err != nil {
		return Keys{}, fmt.Errorf("could not load public key: %w", err)
	}
	return Keys{PrivateKey: privateKey, PublicKey: publicKey}, nil
}

// SHA1Encode computes and returns the SHA-1 hash of the given byte slice.
func SHA1Encode(val []byte) []byte {
	h := sha1.New()
	h.Write(val)
	return h.Sum(nil)
}
