package local

import (
	"bytes"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestAuthenticator(t *testing.T) {
	configBytes := []byte(`
gooduser: "$2a$10$h8bTe02uZIkAa5j1NiuVVOXdUONmch.y151qyK004Hb8EF7rTRq0u"
`)

	testConfig := viper.New()
	testConfig.SetConfigType("yaml")
	err := testConfig.ReadConfig(bytes.NewBuffer(configBytes))
	if err != nil {
		t.Error(err)
	}

	local := &Authenticator{}
	local.Init(testConfig)

	cases := []struct {
		user     string
		password string
		err      string
		valid    bool
	}{
		{"", "", "empty username", false},
		{"", "badpassword", "empty username", false},
		{"baduser", "", "empty password", false},
		{"baduser", "badpassword", "user not found", false},
		{"gooduser", "badpassword", "bad password", false},
		{"gooduser", "goodpassword", "", true},
	}

	for _, c := range cases {
		valid, err := local.Login(c.user, c.password)

		if !c.valid {
			assert.EqualError(t, err, c.err)
		}

		assert.Equal(t, c.valid, valid)
	}
}
