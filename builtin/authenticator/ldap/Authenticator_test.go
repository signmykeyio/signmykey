package ldap

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAuthenticator(t *testing.T) {
	// TODO: Add LDAP mocking
	t.Skip()

	ldap := &Authenticator{
		Address:      "127.0.0.1",
		Port:         636,
		BindUser:     "CN=fakebinduser,OU=Users,DC=test,DC=domain",
		BindPassword: "fakebindpasswd",
		SearchBase:   "OU=Users,DC=test,DC=domain",
		SearchStr:    "(&(objectClass=organizationalPerson)(sAMAccountName=%s))",
		UseTLS:       true,
		TLSVerify:    true,
	}

	valid, err := ldap.Login("fakeuser", "fakepassword")
	if !valid || err != nil {
		t.Logf("%s", err)
		t.Fail()
		return
	}
}

func TestAuthenticatorInit(t *testing.T) {
	cases := []struct {
		config map[string]string
		auth   Authenticator
		err    string
	}{
		{
			map[string]string{},
			Authenticator{},
			"Missing config entries (ldapAddr, ldapPort, ldapTLS, ldapTLSVerify, ldapBindUser, ldapBindPassword, ldapBase, ldapSearch) for Authenticator",
		},
		{
			map[string]string{"ldapAddr": "127.0.0.1"},
			Authenticator{},
			"Missing config entries (ldapPort, ldapTLS, ldapTLSVerify, ldapBindUser, ldapBindPassword, ldapBase, ldapSearch) for Authenticator",
		},
		{
			map[string]string{"ldapAddr": "127.0.0.1", "ldapSearch": "(&(objectClass=organizationalPerson)(sAMAccountName=%s))"},
			Authenticator{},
			"Missing config entries (ldapPort, ldapTLS, ldapTLSVerify, ldapBindUser, ldapBindPassword, ldapBase) for Authenticator",
		},
		{
			map[string]string{
				"ldapAddr":         "127.0.0.1",
				"ldapPort":         "636",
				"ldapTLS":          "True",
				"ldapTLSVerify":    "True",
				"ldapBindUser":     "binduser",
				"ldapBindPassword": "bindpassword",
				"ldapBase":         "DC=fake,DC=org",
				"ldapSearch":       "(&(objectClass=organizationalPerson)(sAMAccountName=%s))",
			},
			Authenticator{
				Address:      "127.0.0.1",
				Port:         636,
				UseTLS:       true,
				TLSVerify:    true,
				BindUser:     "binduser",
				BindPassword: "bindpassword",
				SearchBase:   "DC=fake,DC=org",
				SearchStr:    "(&(objectClass=organizationalPerson)(sAMAccountName=%s))",
			},
			"",
		},
		{
			map[string]string{
				"ldapAddr":         "myldapserver.local",
				"ldapPort":         "389",
				"ldapTLS":          "False",
				"ldapTLSVerify":    "False",
				"ldapBindUser":     "binduser",
				"ldapBindPassword": "bindpassword",
				"ldapBase":         "DC=fake,DC=org",
				"ldapSearch":       "(&(objectClass=organizationalPerson)(sAMAccountName=%s))",
			},
			Authenticator{
				Address:      "myldapserver.local",
				Port:         389,
				UseTLS:       false,
				TLSVerify:    false,
				BindUser:     "binduser",
				BindPassword: "bindpassword",
				SearchBase:   "DC=fake,DC=org",
				SearchStr:    "(&(objectClass=organizationalPerson)(sAMAccountName=%s))",
			},
			"",
		},
		{
			map[string]string{
				"ldapAddr":         "myldapserver.local",
				"ldapPort":         "",
				"ldapTLS":          "False",
				"ldapTLSVerify":    "False",
				"ldapBindUser":     "binduser",
				"ldapBindPassword": "",
				"ldapBase":         "DC=fake,DC=org",
				"ldapSearch":       "(&(objectClass=organizationalPerson)(sAMAccountName=%s))",
			},
			Authenticator{},
			"Missing config entries (ldapPort, ldapBindPassword) for Authenticator",
		},
	}

	for _, c := range cases {
		auth := Authenticator{}
		err := auth.Init(c.config)

		assert.EqualValues(t, c.auth, auth)
		if c.err == "" {
			assert.NoError(t, err)
		} else {
			assert.EqualError(t, err, c.err)
		}
	}
}
