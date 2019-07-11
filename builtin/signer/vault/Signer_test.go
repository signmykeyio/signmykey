package vault

import (
	"context"
	"fmt"
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/ssh"
)

func TestSigner(t *testing.T) {

	// TODO: Add Vault API Call mocking
	t.Skip()

	vs := &Signer{
		Address:  "127.0.0.1",
		Port:     8200,
		RoleID:   "smkroleid",
		SecretID: "smksecretid",
		Path:     "smk",
		Role:     "smk",
	}

	testKey := "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQDCScycMJgKYj7IsyrPYsCVryhz4//mjekvElmihYLc/njL1cC9KBTRhdbV1NYw9RFC/CENAwrHcGXAcgMuY0fFUQOyKDa6HmVIxT8vszcePAutl6YvcFYuTCJRYjOQXWBmYEe1NJI8yfR+CMU8HdfXbXUhU93UxpNX8rzXGv3KSSky6v1BAkzTl5QiQRdjn18Wf8z3e1sMbxjfD+ygas9rpmbUONsUc6d5U6YYTroJA/DoRIG92IBK9GjH1w/l9Wvcs2V5atCorKwasdCsEFO5Jv84XO41smo/IF+gd0hUtKDDGkk/djk3TmC9h2WUBb43lxiv0wn2ByTIMAorCqFb" // nolint: lll

	cases := []struct {
		payload    []byte
		id         string
		principals []string
		expErr     bool
	}{
		{[]byte("{\"public_key\": \"invalid key\"}"), "test", []string{"root", "admin"}, true},
		{[]byte(fmt.Sprintf("{\"public_key\": \"%s\"}", testKey)), "testid", []string{"admin", "root"}, false},
		{[]byte(fmt.Sprintf("{\"public_key\": \"%s\"}", testKey)), "testid", []string{"root", "admin"}, false},
		{[]byte(""), "testid", []string{"root", "admin"}, true},
		{[]byte(fmt.Sprintf("{\"public_key\": \"%s\"}", testKey)), "", []string{"root", "admin"}, true},
		{[]byte(fmt.Sprintf("{\"public_key\": \"%s\"}", testKey)), "testid", []string{}, true},
	}

	for _, c := range cases {
		cert, err := vs.Sign(context.Background(), c.payload, c.id, c.principals)
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

		sort.Strings(c.principals)
		sort.Strings(sshCert.ValidPrincipals)

		assert.Equal(t, c.id, sshCert.KeyId)
		assert.Equal(t, c.principals, sshCert.ValidPrincipals)
	}
}
