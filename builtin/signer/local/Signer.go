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
			return fmt.Errorf("Config entry %s missing for Signer", entry)
		}
	}

	// Read and parse CA private key
	key, err := ioutil.ReadFile(config.GetString("caKey"))
	if err != nil {
		return err
	}
	s.CAKey, err = ssh.ParsePrivateKey(key)
	if err != nil {
		return err
	}

	// Read and parse CA public key
	pubKey, err := ioutil.ReadFile(config.GetString("caCert"))
	if err != nil {
		return err
	}
	s.CACert, err = ssh.ParsePublicKey(pubKey)
	if err != nil {
		return err
	}

	s.TTL = viper.GetInt("ttl")
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

	err := checkCertReqFields(certreq)
	if err != nil {
		return "", errors.Wrap(err, "certificate request fields checking failed")
	}

	buf := make([]byte, 8)
	rand.Read(buf)
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
		ValidAfter:      uint64(time.Now().Add(-1 * time.Minute).In(time.UTC).Unix()),
		ValidBefore:     uint64(time.Now().Add(time.Duration(s.TTL) * time.Second).In(time.UTC).Unix()),
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

func checkCertReqFields(certreq signer.CertReq) error {
	if len(certreq.Principals) == 0 {
		return fmt.Errorf("Empty list of principals")
	}
	if len(certreq.ID) == 0 {
		return fmt.Errorf("Empty ID string")
	}
	if len(certreq.Key) == 0 {
		return fmt.Errorf("Empty key string")
	}

	return nil
}
