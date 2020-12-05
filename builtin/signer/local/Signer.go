package local

import (
	"context"
	"crypto/rand"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"time"

	"github.com/signmykeyio/signmykey/builtin/signer"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"golang.org/x/crypto/ssh"
)

// Signer struct represents local options for signing SSH Key.
type Signer struct {
	CACert          ssh.PublicKey
	CAKey           ssh.Signer
	TTL             int
	CriticalOptions map[string]string
	Extensions      map[string]string
}

type localSignReq struct {
	PubKey string `json:"public_key" binding:"required"`
}

// Init method is used to ingest config of Signer
func (s *Signer) Init(config *viper.Viper) error {
	neededEntries := []string{
		"caCert",
		"caKey",
		"ttl",
	}

	for _, entry := range neededEntries {
		if !config.IsSet(entry) {
			return fmt.Errorf("config entry %s missing for Signer", entry)
		}
	}

	// Read and parse CA private key
	key, err := ioutil.ReadFile(config.GetString("caKey"))
	if err != nil {
		return fmt.Errorf("error reading CA private key file %s: %w", config.GetString("caKey"), err)
	}
	s.CAKey, err = ssh.ParsePrivateKey(key)
	if err != nil {
		return fmt.Errorf("error parsing CA private key: %w", err)
	}

	// Read and parse CA public key
	pubKey, err := ioutil.ReadFile(config.GetString("caCert"))
	if err != nil {
		return fmt.Errorf("error reading CA public key file %s: %w", config.GetString("caCert"), err)
	}
	s.CACert, _, _, _, err = ssh.ParseAuthorizedKey(pubKey)
	if err != nil {
		return fmt.Errorf("error parsing CA public key: %w", err)
	}

	config.SetDefault("extensions", map[string]string{
		"permit-X11-forwarding":   "",
		"permit-agent-forwarding": "",
		"permit-port-forwarding":  "",
		"permit-pty":              "",
		"permit-user-rc":          "",
	})

	s.TTL = config.GetInt("ttl")
	s.CriticalOptions = config.GetStringMapString("criticalOptions")
	s.Extensions = config.GetStringMapString("extensions")

	return nil
}

// ReadCA method read CA public cert from local file
func (s Signer) ReadCA() (string, error) {
	return string(ssh.MarshalAuthorizedKey(s.CACert)), nil
}

// Sign method is used to sign passed SSH Key.
func (s Signer) Sign(ctx context.Context, payload []byte, id string, principals []string) (cert string, err error) {

	var signReq localSignReq
	err = json.Unmarshal(payload, &signReq)
	if err != nil {
		log.Errorf("json unmarshaling failed: %s", err)
		return "", fmt.Errorf("JSON unmarshaling failed: %w", err)
	}

	if id == "" {
		return "", errors.New("empty id")
	}

	if len(principals) == 0 {
		return "", errors.New("empty list of principals")
	}

	certreq := signer.CertReq{
		Key:        signReq.PubKey,
		ID:         id,
		Principals: principals,
	}
	buf := make([]byte, 8)
	_, err = rand.Read(buf)
	if err != nil {
		return "", fmt.Errorf("failed to read random bytes: %w", err)
	}
	serial := binary.LittleEndian.Uint64(buf)

	pubKey, _, _, _, err := ssh.ParseAuthorizedKey([]byte(certreq.Key))
	if err != nil {
		return "", fmt.Errorf("failed to parse user public key: %w", err)
	}

	certificate := ssh.Certificate{
		Serial:          serial,
		Key:             pubKey,
		KeyId:           certreq.ID,
		ValidPrincipals: certreq.Principals,
		ValidAfter:      uint64(time.Now().Unix() - 60),
		ValidBefore:     uint64(time.Now().Unix() + int64(s.TTL)),
		CertType:        ssh.UserCert,
		Permissions: ssh.Permissions{
			CriticalOptions: s.CriticalOptions,
			Extensions:      s.Extensions,
		},
	}

	err = certificate.SignCert(rand.Reader, s.CAKey)
	if err != nil {
		return "", fmt.Errorf("failed to sign user public key: %w", err)
	}

	marshaledCertificate := ssh.MarshalAuthorizedKey(&certificate)
	if len(marshaledCertificate) == 0 {
		return "", errors.New("failed to marshal signed certificate, empty result")
	}

	return string(marshaledCertificate), nil
}
