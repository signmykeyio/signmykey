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
			"Missing config entries (ldap-addr, ldap-port, ldap-tls, ldap-tls-verify, ldap-bind-user, ldap-bind-password, ldap-base, ldap-search) for Authenticator",
		},
		{
			map[string]string{"ldap-addr": "127.0.0.1"},
			Authenticator{},
			"Missing config entries (ldap-port, ldap-tls, ldap-tls-verify, ldap-bind-user, ldap-bind-password, ldap-base, ldap-search) for Authenticator",
		},
		{
			map[string]string{"ldap-addr": "127.0.0.1", "ldap-search": "(&(objectClass=organizationalPerson)(sAMAccountName=%s))"},
			Authenticator{},
			"Missing config entries (ldap-port, ldap-tls, ldap-tls-verify, ldap-bind-user, ldap-bind-password, ldap-base) for Authenticator",
		},
		{
			map[string]string{
				"ldap-addr":          "127.0.0.1",
				"ldap-port":          "636",
				"ldap-tls":           "True",
				"ldap-tls-verify":    "True",
				"ldap-bind-user":     "binduser",
				"ldap-bind-password": "bindpassword",
				"ldap-base":          "DC=fake,DC=org",
				"ldap-search":        "(&(objectClass=organizationalPerson)(sAMAccountName=%s))",
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
				"ldap-addr":          "myldapserver.local",
				"ldap-port":          "389",
				"ldap-tls":           "False",
				"ldap-tls-verify":    "False",
				"ldap-bind-user":     "binduser",
				"ldap-bind-password": "bindpassword",
				"ldap-base":          "DC=fake,DC=org",
				"ldap-search":        "(&(objectClass=organizationalPerson)(sAMAccountName=%s))",
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
				"ldap-addr":          "myldapserver.local",
				"ldap-port":          "",
				"ldap-tls":           "False",
				"ldap-tls-verify":    "False",
				"ldap-bind-user":     "binduser",
				"ldap-bind-password": "",
				"ldap-base":          "DC=fake,DC=org",
				"ldap-search":        "(&(objectClass=organizationalPerson)(sAMAccountName=%s))",
			},
			Authenticator{},
			"Missing config entries (ldap-port, ldap-bind-password) for Authenticator",
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
