package client

import (
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

	pubKey, err := ioutil.ReadFile(pubKeyPath) // nolint: gosec
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

		return fmt.Errorf(`user SSH key at %s doesn't exist.
Please generate one with this command :

    ssh-keygen -f %s`,
			pubKeyPath, keyPath)
	}

	return nil
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

	cert, err := ioutil.ReadFile(fullPath) // nolint: gosec
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

// CertInfo extract principals and expiration from SSH certificate
func CertInfo(cert string) (principals []string, before uint64, err error) {
	parsedKey, _, _, _, err := ssh.ParseAuthorizedKey([]byte(cert))
	if err != nil {
		return principals, before, err
	}

	parsedCert := parsedKey.(*ssh.Certificate)

	return parsedCert.ValidPrincipals, parsedCert.ValidBefore, nil
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

	err = os.Chmod(signedKeyPath, 0644)
	if err != nil {
		return err
	}

	_, err = signedKeyFile.WriteString(signedKey)
	return err
}
