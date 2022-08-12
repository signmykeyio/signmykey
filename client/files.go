package client

import (
	"fmt"
	"os"
	"reflect"
	"strings"
	"time"

	"golang.org/x/crypto/ssh"

	homedir "github.com/mitchellh/go-homedir"
)

// DefaultSSHKeys is based on `man ssh` -i identity_file default values
var DefaultSSHKeys = []string{
	"~/.ssh/id_dsa.pub",
	"~/.ssh/id_ecdsa.pub",
	"~/.ssh/id_ecdsa_sk.pub",
	"~/.ssh/id_ed25519.pub",
	"~/.ssh/id_ed25519_sk.pub",
	"~/.ssh/id_rsa.pub",
}

// GetUserPubKey returns user's SSH public key as string.
func GetUserPubKey(key string) (string, error) {
	pubKeyPath, err := homedir.Expand(key)
	if err != nil {
		return "", err
	}

	pubKey, err := os.ReadFile(pubKeyPath) // nolint: gosec
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(string(pubKey)), nil
}

// UserPubKeyExists checks if public key file exists.
func UserPubKeyExists(key string) (string, error) {
	pubKeyPath, err := homedir.Expand(key)
	if err != nil {
		return "", err
	}

	if _, err := os.Stat(pubKeyPath); os.IsNotExist(err) {
		return "", nil
	}

	return pubKeyPath, nil
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

	cert, err := os.ReadFile(fullPath) // nolint: gosec
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

	err = os.Chmod(signedKeyPath, 0600)
	if err != nil {
		return err
	}

	_, err = signedKeyFile.WriteString(signedKey)
	return err
}

// FindUserPubKeys checks every pubkey in `keys` list and returns only existsing keys (or error if
// all pubkeys don't exist)
func FindUserPubKeys(keys []string) ([]string, error) {
	var found = []string{}
	for _, key := range keys {
		pubKeyPath, err := UserPubKeyExists(key)
		if err != nil {
			return nil, err
		}
		if pubKeyPath != "" {
			found = append(found, pubKeyPath)
		}
	}

	if len(found) == 0 {
		var errStr = fmt.Sprintf(`user SSH keys at %s doesn't exist.

Please generate at least one with command like this :

`, strings.Join(keys, ", "))

		// keys list is not explicitly set, suggest generating ed25519 as default
		if reflect.DeepEqual(keys, DefaultSSHKeys) {
			errStr += "\tssh-keygen -f ~/.ssh/id_ed25519 -t ed25519\n"
		} else {
			for _, key := range keys {
				suggestedSSHKey := strings.Replace(key, ".pub", "", 1)
				errStr += fmt.Sprintf("\tssh-keygen -f %s -t %s\n", suggestedSSHKey, chooseSSHKeyType(suggestedSSHKey))
			}
		}

		return nil, fmt.Errorf(errStr)
	}
	return found, nil
}

func chooseSSHKeyType(key string) string {
	switch {
	case strings.Contains(key, "ecdsa_sk"):
		return "ecdsa-sk"
	case strings.Contains(key, "ed25519_sk"):
		return "ed25519-sk"
	case strings.Contains(key, "ecdsa"):
		return "ecdsa"
	case strings.Contains(key, "dsa"):
		return "dsa"
	case strings.Contains(key, "rsa"):
		return "rsa"
	default:
		return "ed25519"
	}
}
