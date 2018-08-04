package local

import (
	"crypto/rand"
	"encoding/binary"
	"fmt"
	"io/ioutil"
	"strconv"
	"time"

	"github.com/pkg/errors"
	"github.com/signmykeyio/signmykey/builtin/signer"
	"golang.org/x/crypto/ssh"
)

// Signer struct represents Hashicorp Vault options for signing SSH Key.
type Signer struct {
	CACert          string
	CAKey           string
	TTL             int
	CriticalOptions map[string]string
	Extensions      map[string]string
}

// Init method is used to ingest config of Signer
func (s *Signer) Init(config map[string]string) error {
	neededEntries := []string{
		"caCert",
		"caKey",
		"ttl",
	}

	for _, entry := range neededEntries {
		if _, ok := config[entry]; !ok {
			return fmt.Errorf("Config entry %s missing for Signer", entry)
		}
	}

	// Conversions
	ttl, err := strconv.Atoi(config["ttl"])
	if err != nil {
		return err
	}

	s.CACert = config["caCert"]
	s.CAKey = config["caKey"]
	s.TTL = ttl
	s.CriticalOptions = map[string]string{}
	s.Extensions = map[string]string{}

	return nil
}

// ReadCA method read CA public cert from local file
func (s Signer) ReadCA() (string, error) {
	publicKey, err := ioutil.ReadFile(s.CACert)

	return string(publicKey), err
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

	pubKey, err := ssh.ParsePublicKey([]byte(certreq.Key))
	if err != nil {
		return "", err
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
			CriticalOptions: map[string]string{},
			Extensions:      map[string]string{},
		},
	}

	ca, err := ioutil.ReadFile(s.CAKey)
	if err != nil {
		return "", err
	}

	signer, err := ssh.ParsePrivateKey(ca)
	if err != nil {
		return "", err
	}

	err = certificate.SignCert(rand.Reader, signer)
	if err != nil {
		return "", err
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
