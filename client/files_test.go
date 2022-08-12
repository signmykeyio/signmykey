package client

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestChooseSSHKeyType(t *testing.T) {
	cases := []struct {
		keyName       string
		keyType       string
		keyDeprecated bool
	}{
		{"~/.ssh/id_dsa.pub", "dsa", true},
		{"~/.ssh/id_ecdsa.pub", "ecdsa", false},
		{"~/.ssh/id_ecdsa_sk.pub", "ecdsa-sk", false},
		{"~/.ssh/id_ed25519.pub", "ed25519", false},
		{"~/.ssh/id_ed25519_sk.pub", "ed25519-sk", false},
		{"~/.ssh/id_rsa.pub", "rsa", false},
		{"~/.ssh/test_default_type.pub", "ed25519", false},
	}

	for _, c := range cases {
		keyType, isDeprecated := chooseSSHKeyType(c.keyName)
		assert.Equal(t, c.keyType, keyType)
		assert.Equal(t, c.keyDeprecated, isDeprecated)
	}
}
