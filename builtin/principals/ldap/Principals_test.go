package ldap

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPrincipals(t *testing.T) {
	// TODO: Add LDAP mocking
	t.Skip()

	ldap := &Principals{
		Address:      "127.0.0.1",
		Port:         636,
		BindUser:     "CN=fakebinduser,OU=Users,DC=test,DC=domain",
		BindPassword: "fakebindpasswd",
		SearchBase:   "OU=Users,DC=test,DC=domain",
		SearchStr:    "(&(objectClass=organizationalPerson)(sAMAccountName=%s))",
		UseTLS:       true,
		TLSVerify:    true,
		Prefix:       "smk-",
	}

	principals, err := ldap.Get("fakeuser")
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

func TestGetCN(t *testing.T) {
	cases := []struct {
		list    []string
		expList []string
	}{
		{
			[]string{
				"CN=grouptest1,OU=Groups,DC=test,DC=domain",
				"CN=grouptest-2,OU=Groups,DC=test,DC=domain",
				"DN=group3,OU=Groups,DC=test,DC=domain",
				"CN=,OU=Groups,DC=test,DC=domain",
				"CN=group4_test,CN=Groups,DC=test,DC=domain",
			},
			[]string{"grouptest1", "grouptest-2", "group4_test"},
		},
	}

	for _, c := range cases {
		cnList := getCN(c.list)
		assert.Equal(t, c.expList, cnList)
	}
}

func TestFilterByPrefix(t *testing.T) {
	cases := []struct {
		prefix  string
		list    []string
		expList []string
	}{
		{"smk-",
			[]string{"group1", "smk-group2", "group3-smk", "smk-group4", ""},
			[]string{"group2", "group4"},
		},
		{"",
			[]string{"group1", "smk-group2", "group3-smk", "smk-group4", ""},
			[]string{"group1", "smk-group2", "group3-smk", "smk-group4", ""},
		},
		{"smk-",
			[]string{},
			[]string{},
		},
	}

	for _, c := range cases {
		filtList := filterByPrefix(c.prefix, c.list)
		assert.Equal(t, c.expList, filtList)
	}
}

func TestPrincipalsInit(t *testing.T) {
	cases := []struct {
		config map[string]string
		auth   Principals
		err    string
	}{
		{
			map[string]string{},
			Principals{},
			"Missing config entries (ldap-addr, ldap-port, ldap-tls, ldap-tls-verify, ldap-bind-user, ldap-bind-password, ldap-base, ldap-search) for Principals",
		},
		{
			map[string]string{"ldap-addr": "127.0.0.1"},
			Principals{},
			"Missing config entries (ldap-port, ldap-tls, ldap-tls-verify, ldap-bind-user, ldap-bind-password, ldap-base, ldap-search) for Principals",
		},
		{
			map[string]string{"ldap-addr": "127.0.0.1", "ldap-search": "(&(objectClass=organizationalPerson)(sAMAccountName=%s))"},
			Principals{},
			"Missing config entries (ldap-port, ldap-tls, ldap-tls-verify, ldap-bind-user, ldap-bind-password, ldap-base) for Principals",
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
			Principals{
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
			Principals{
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
			Principals{},
			"Missing config entries (ldap-port, ldap-bind-password) for Principals",
		},
	}

	for _, c := range cases {
		auth := Principals{}
		err := auth.Init(c.config)

		assert.EqualValues(t, c.auth, auth)
		if c.err == "" {
			assert.NoError(t, err)
		} else {
			assert.EqualError(t, err, c.err)
		}
	}

}
