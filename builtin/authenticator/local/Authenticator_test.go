package local

import (
	"bytes"
	"context"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestAuthenticator(t *testing.T) {
	configBytes := []byte(`
users:
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
		payload []byte
		id      string
		err     string
		valid   bool
	}{
		{[]byte(""), "", "JSON unmarshaling failed: unexpected end of JSON input", false},
		{[]byte("{\"user\":\"\"}"), "", "empty username", false},
		{[]byte("{\"password\":\"\"}"), "", "empty username", false},
		{[]byte("{\"user\":\"baduser\"}"), "", "empty password", false},
		{[]byte("{\"user\":\"baduser\",\"password\":\"badpassword\"}"), "", "user not found", false},
		{[]byte("{\"user\":\"gooduser\",\"password\":\"badpassword\"}"), "", "bad password", false},
		{[]byte("{\"user\":\"gooduser\",\"password\":\"goodpassword\"}"), "local-gooduser", "", true},
	}

	for _, c := range cases {
		_, valid, id, err := local.Login(context.Background(), c.payload)

		assert.Equal(t, c.valid, valid)

		assert.Equal(t, id, c.id)

		if !c.valid {
			assert.EqualError(t, err, c.err)
		}
	}
}
