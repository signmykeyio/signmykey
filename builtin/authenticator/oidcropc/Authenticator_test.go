package oidcropc

import (
	"bytes"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestAuthenticator(t *testing.T) {
	// TODO: Add OIDC ROPC mocking
	t.Skip()

	oidc := &Authenticator{
		OIDCTokenEndpoint: "https://127.0.0.1/auth/realms/master/protocol/openid-connect/token",
		OIDCClientID:      "signmykey",
		OIDCClientSecret:  "461c5204-160b-45cf-8609-aa5d500e6093",
	}

	valid, _, err := oidc.Login("fakeuser", "fakepassword")
	if !valid || err != nil {
		t.Logf("%s", err)
		t.Fail()
		return
	}
}

func TestAuthenticatorInit(t *testing.T) {
	cases := []struct {
		config []byte
		auth   Authenticator
		err    string
	}{
		{
			[]byte(""),
			Authenticator{},
			"Missing config entries (oidcTokenEndpoint, oidcClientID, oidcClientSecret) for Authenticator",
		},
		{
			[]byte(`oidcTokenEndpoint: "https://127.0.0.1/auth/realms/master/protocol/openid-connect/token"`),
			Authenticator{},
			"Missing config entries (oidcClientID, oidcClientSecret) for Authenticator",
		},
		{
			[]byte(`
oidcTokenEndpoint: "https://127.0.0.1/auth/realms/master/protocol/openid-connect/token"
oidcClientID: "signmykey"
`),
			Authenticator{},
			"Missing config entries (oidcClientSecret) for Authenticator",
		},
		{
			[]byte(`
oidcTokenEndpoint: "https://127.0.0.1/auth/realms/master/protocol/openid-connect/token"
oidcClientID: "signmykey"
oidcClientSecret: "461c5204-160b-45cf-8609-aa5d500e6093"
`),
			Authenticator{
				OIDCTokenEndpoint: "https://127.0.0.1/auth/realms/master/protocol/openid-connect/token",
				OIDCClientID:      "signmykey",
				OIDCClientSecret:  "461c5204-160b-45cf-8609-aa5d500e6093",
			},
			"",
		},
	}

	for _, c := range cases {
		testConfig := viper.New()
		testConfig.SetConfigType("yaml")
		err := testConfig.ReadConfig(bytes.NewBuffer(c.config))
		if err != nil {
			t.Error(err)
		}

		auth := Authenticator{}
		err = auth.Init(testConfig)

		assert.EqualValues(t, c.auth, auth)
		if c.err == "" {
			assert.NoError(t, err)
		} else {
			assert.EqualError(t, err, c.err)
		}
	}
}
