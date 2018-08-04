package local

import (
	"sort"
	"testing"

	"github.com/signmykeyio/signmykey/builtin/signer"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/ssh"
)

func TestSigner(t *testing.T) {

	// TODO: mock CA files
	t.Skip()

	s := &Signer{
		CACert: "ca",
		CAKey:  "ca.key",
		TTL:    600,
	}

	testKey := "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQDCScycMJgKYj7IsyrPYsCVryhz4//mjekvElmihYLc/njL1cC9KBTRhdbV1NYw9RFC/CENAwrHcGXAcgMuY0fFUQOyKDa6HmVIxT8vszcePAutl6YvcFYuTCJRYjOQXWBmYEe1NJI8yfR+CMU8HdfXbXUhU93UxpNX8rzXGv3KSSky6v1BAkzTl5QiQRdjn18Wf8z3e1sMbxjfD+ygas9rpmbUONsUc6d5U6YYTroJA/DoRIG92IBK9GjH1w/l9Wvcs2V5atCorKwasdCsEFO5Jv84XO41smo/IF+gd0hUtKDDGkk/djk3TmC9h2WUBb43lxiv0wn2ByTIMAorCqFb" // nolint: lll

	cases := []struct {
		req    signer.CertReq
		expErr bool
	}{
		{signer.CertReq{"invalid key", "test", []string{"root", "admin"}}, true},
		{signer.CertReq{testKey, "testid", []string{"admin", "root"}}, false},
		{signer.CertReq{testKey, "testid", []string{"root", "admin"}}, false},
		{signer.CertReq{"", "testid", []string{"root", "admin"}}, true},
		{signer.CertReq{testKey, "", []string{"root", "admin"}}, true},
		{signer.CertReq{testKey, "testid", []string{}}, true},
	}

	for _, c := range cases {
		cert, err := s.Sign(c.req)
		if c.expErr {
			assert.Error(t, err)
		} else {
			assert.NoError(t, err)
		}

		if err != nil {
			continue
		}

		parsedCert, _, _, _, _ := ssh.ParseAuthorizedKey([]byte(cert))
		sshCert := parsedCert.(*ssh.Certificate)

		sort.Strings(c.req.Principals)
		sort.Strings(sshCert.ValidPrincipals)

		assert.Equal(t, c.req.ID, sshCert.KeyId)
		assert.Equal(t, c.req.Principals, sshCert.ValidPrincipals)
	}
}
