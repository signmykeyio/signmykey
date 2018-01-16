package helper

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"time"

	"golang.org/x/crypto/ssh"

	homedir "github.com/mitchellh/go-homedir"
)

// GetUserPubKey returns user's SSH public key as string.
func GetUserPubKey(key string) (string, error) {
	pubKeyPath, err := homedir.Expand(key)
	if err != nil {
		return "", err
	}

	pubKey, err := ioutil.ReadFile(pubKeyPath)
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(string(pubKey)), nil
}

// UserPubKeyExists checks if public key file exists.
func UserPubKeyExists(key string) error {
	pubKeyPath, err := homedir.Expand(key)
	if err != nil {
		return err
	}

	if _, err := os.Stat(pubKeyPath); os.IsNotExist(err) {
		keyPath := strings.Replace(pubKeyPath, ".pub", "", 1)
		err := generateSSHKeyPair(keyPath)
		if err != nil {
			return err
		}
	}

	return nil
}

func generateSSHKeyPair(path string) error {
	// Generate RSA private key
	private, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return err
	}

	// Generate SSH public key from RSA private key
	public, err := ssh.NewPublicKey(&private.PublicKey)
	if err != nil {
		return err
	}

	// Convert RSA private key to ssh format and write it
	privateBlock := &pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(private)}
	privatePEM := pem.EncodeToMemory(privateBlock)
	err = ioutil.WriteFile(path, privatePEM, 0600)
	if err != nil {
		return err
	}

	// Convert and write SSH public key
	serializedPublicKey := ssh.MarshalAuthorizedKey(public)
	err = ioutil.WriteFile(fmt.Sprintf("%s.pub", path), serializedPublicKey, 0644)
	return err
}

// CertStillValid checks if the certificate is not expired.
func CertStillValid(path string) bool {
	fullPath, err := homedir.Expand(strings.Replace(path, ".pub", "-cert.pub", 1))
	if err != nil {
		return false
	}

	if _, err = os.Stat(fullPath); os.IsNotExist(err) {
		return false
	}

	cert, err := ioutil.ReadFile(fullPath)
	if err != nil {
		return false
	}

	parsedKey, _, _, _, err := ssh.ParseAuthorizedKey(cert)
	if err != nil {
		return false
	}

	parsedCert := parsedKey.(*ssh.Certificate)

	return parsedCert.ValidBefore > uint64(time.Now().Unix())
}

// WriteUserSignedKey writes user certificate on disk.
func WriteUserSignedKey(signedKey string, key string) (err error) {
	signedKeyPath, err := homedir.Expand(strings.Replace(key, ".pub", "-cert.pub", 1))
	if err != nil {
		return err
	}

	signedKeyFile, err := os.Create(signedKeyPath)
	if err != nil {
		return err
	}
	defer func() {
		if cerr := signedKeyFile.Close(); cerr != nil && err == nil {
			err = cerr
		}
	}()

	err = signedKeyFile.Chmod(0644)
	if err != nil {
		return err
	}

	_, err = signedKeyFile.WriteString(signedKey)
	return err
}
