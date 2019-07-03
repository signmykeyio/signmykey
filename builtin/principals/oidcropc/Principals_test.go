package oidcropc

import (
	"bytes"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestPrincipals(t *testing.T) {
	// TODO: Add OIDC ROPC mocking
	t.Skip()

	oidc := &Principals{
		OIDCUserinfoEndpoint: "https://127.0.0.1/auth/realms/master/protocol/openid-connect/userinfo",
		OIDCUserGroupsEntry:  "oidc-groups",
	}

	principals, err := oidc.Get("faketoken")
	if err != nil {
		t.Logf("%s", err)
		t.Fail()
		return
	} else if len(principals) == 0 {
		t.Logf("empty list of principals")
		t.Fail()
		return
	}
}

func TestPrincipalsInit(t *testing.T) {
	cases := []struct {
		config []byte
		auth   Principals
		err    string
	}{
		{
			[]byte(""),
			Principals{},
			"Missing config entries (oidcUserinfoEndpoint, oidcUserGroupsEntry) for Principals",
		},
		{
			[]byte(`oidcUserinfoEndpoint: "https://127.0.0.1/auth/realms/master/protocol/openid-connect/userinfo"`),
			Principals{},
			"Missing config entries (oidcUserGroupsEntry) for Principals",
		},
		{
			[]byte(`
oidcUserinfoEndpoint: "https://127.0.0.1/auth/realms/master/protocol/openid-connect/userinfo"
oidcUserGroupsEntry: "oidc-groups"
`),
			Principals{
				OIDCUserinfoEndpoint: "https://127.0.0.1/auth/realms/master/protocol/openid-connect/userinfo",
				OIDCUserGroupsEntry:  "oidc-groups",
				TransformCase:        "none",
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

		auth := Principals{}
		err = auth.Init(testConfig)

		assert.EqualValues(t, c.auth, auth)
		if c.err == "" {
			assert.NoError(t, err)
		} else {
			assert.EqualError(t, err, c.err)
		}
	}

}
