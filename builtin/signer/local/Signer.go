package local

import (
	"crypto/rand"
	"encoding/binary"
	"fmt"
	"io/ioutil"
	"time"

	"github.com/pkg/errors"
	"github.com/signmykeyio/signmykey/builtin/signer"
	"github.com/spf13/viper"
	"golang.org/x/crypto/ssh"
)

// Signer struct represents Hashicorp Vault options for signing SSH Key.
type Signer struct {
	CACert          ssh.PublicKey
	CAKey           ssh.Signer
	TTL             int
	CriticalOptions map[string]string
	Extensions      map[string]string
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
		return errors.Wrapf(err, "error reading CA private key file %s", config.GetString("caKey"))
	}
	s.CAKey, err = ssh.ParsePrivateKey(key)
	if err != nil {
		return errors.Wrap(err, "error parsing CA private key")
	}

	// Read and parse CA public key
	pubKey, err := ioutil.ReadFile(config.GetString("caCert"))
	if err != nil {
		return errors.Wrapf(err, "error reading CA public key file %s", config.GetString("caCert"))
	}
	s.CACert, _, _, _, err = ssh.ParseAuthorizedKey(pubKey)
	if err != nil {
		return errors.Wrap(err, "error parsing CA public key")
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
	return string(s.CACert.Marshal()), nil
}

// Sign method is used to sign passed SSH Key.
func (s Signer) Sign(certreq signer.CertReq) (string, error) {

	err := certreq.CheckCertReqFields()
	if err != nil {
		return "", errors.Wrap(err, "certificate request fields checking failed")
	}

	buf := make([]byte, 8)
	_, err = rand.Read(buf)
	if err != nil {
		return "", errors.Wrap(err, "failed to read random bytes")
	}
	serial := binary.LittleEndian.Uint64(buf)

	pubKey, _, _, _, err := ssh.ParseAuthorizedKey([]byte(certreq.Key))
	if err != nil {
		return "", fmt.Errorf("failed to parse user public key: %s", err)
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
		return "", fmt.Errorf("failed to sign user public key: %s", err)
	}

	marshaledCertificate := ssh.MarshalAuthorizedKey(&certificate)
	if len(marshaledCertificate) == 0 {
		return "", errors.New("failed to marshal signed certificate, empty result")
	}

	return string(marshaledCertificate), nil
}
