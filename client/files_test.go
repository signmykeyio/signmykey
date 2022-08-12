package client

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestChooseSSHKeyType(t *testing.T) {
	cases := []struct {
		keyName string
		keyType string
	}{
		{"~/.ssh/id_dsa.pub", "dsa"},
		{"~/.ssh/id_ecdsa.pub", "ecdsa"},
		{"~/.ssh/id_ecdsa_sk.pub", "ecdsa-sk"},
		{"~/.ssh/id_ed25519.pub", "ed25519"},
		{"~/.ssh/id_ed25519_sk.pub", "ed25519-sk"},
		{"~/.ssh/id_rsa.pub", "rsa-sha2-512"},
		{"~/.ssh/test_default_type.pub", "ed25519"},
	}

	for _, c := range cases {
		keyType := chooseSSHKeyType(c.keyName)
		assert.Equal(t, c.keyType, keyType)
	}
}
