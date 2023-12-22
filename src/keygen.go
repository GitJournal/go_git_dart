package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"

	"errors"

	"golang.org/x/crypto/ssh"
)

// generateRSAKeys generates an RSA public/private key pair
// and returns them as PEM encoded strings.
func generateRSAKeys() (string, string, error) {
	// Generate the private key
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048*2)
	if err != nil {
		return "", "", err
	}

	// Marshal the private key to DER format
	privateKeyDER := x509.MarshalPKCS1PrivateKey(privateKey)

	// PEM encode the private key
	privateKeyBlock := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: privateKeyDER,
	}
	privateKeyPEM := string(pem.EncodeToMemory(privateKeyBlock))

	// Generate the public key
	publicKey := &privateKey.PublicKey

	// Marshal the public key to DER format
	publicKeyDER, err := x509.MarshalPKIXPublicKey(publicKey)
	if err != nil {
		return "", "", err
	}

	// PEM encode the public key
	publicKeyBlock := &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: publicKeyDER,
	}
	publicKeyPEM := string(pem.EncodeToMemory(publicKeyBlock))
	publicKeyOpenSSH, err := PublicPEMtoOpenSSH([]byte(publicKeyPEM))
	if err != nil {
		return "", "", err
	}

	return publicKeyOpenSSH, privateKeyPEM, nil
}

// Converts PEM public key to OpenSSH format to be used in authorized_keys file
// Similar to: "ssh-keygen", "-i", "-m", "pkcs8", "-f", auth_keys_new_path
func PublicPEMtoOpenSSH(pemBytes []byte) (string, error) {
	// Decode and get the first block in the PEM file.
	// In our case it should be the Public key block.
	pemBlock, rest := pem.Decode(pemBytes)
	if pemBlock == nil {
		return "", errors.New("invalid PEM public key passed, pem.Decode() did not find a public key")
	}
	if len(rest) > 0 {
		return "", errors.New("PEM block contains more than just public key")
	}

	// Confirm we got the PUBLIC KEY block type
	if pemBlock.Type != "PUBLIC KEY" {
		return "", fmt.Errorf("ssh: unsupported key type %q", pemBlock.Type)
	}

	// Convert to rsa
	rsaPubKey, err := x509.ParsePKIXPublicKey(pemBlock.Bytes)
	if err != nil {
		return "", fmt.Errorf("x509.parse pki public key: %w", err)
	}

	// Confirm we got an rsa public key. Returned value is an interface{}
	sshKey, ok := rsaPubKey.(*rsa.PublicKey)
	if !ok {
		return "", fmt.Errorf("invalid PEM passed in from user: %w", err)
	}

	// Generate the ssh public key
	pub, err := ssh.NewPublicKey(sshKey)
	if err != nil {
		return "", fmt.Errorf("new ssh public key from rsa: %w", err)
	}

	// Encode to store to file
	sshPubKey := base64.StdEncoding.EncodeToString(pub.Marshal())

	return "ssh-rsa " + sshPubKey, nil
}
