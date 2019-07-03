package ldap

import (
	"crypto/tls"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/spf13/viper"
	ldap "gopkg.in/ldap.v2"
)

// Authenticator struct represents LDAP options for SMK Authentication.
type Authenticator struct {
	Address      string
	Port         int
	BindUser     string
	BindPassword string
	SearchBase   string
	SearchStr    string
	UseTLS       bool
	TLSVerify    bool
}

// Init method is used to ingest config of Authenticator
func (a *Authenticator) Init(config *viper.Viper) error {
	neededEntries := []string{
		"ldapAddr",
		"ldapPort",
		"ldapTLS",
		"ldapTLSVerify",
		"ldapBindUser",
		"ldapBindPassword",
		"ldapBase",
		"ldapSearch",
	}

	var missingEntriesLst []string
	for _, entry := range neededEntries {
		if !config.IsSet(entry) {
			missingEntriesLst = append(missingEntriesLst, entry)
		}
	}
	if len(missingEntriesLst) > 0 {
		missingEntries := strings.Join(missingEntriesLst, ", ")
		return fmt.Errorf("Missing config entries (%s) for Authenticator", missingEntries)
	}

	a.Address = config.GetString("ldapAddr")
	a.Port = config.GetInt("ldapPort")
	a.UseTLS = config.GetBool("ldapTLS")
	a.TLSVerify = config.GetBool("ldapTLSVerify")
	a.BindUser = config.GetString("ldapBindUser")
	a.BindPassword = config.GetString("ldapBindPassword")
	a.SearchBase = config.GetString("ldapBase")
	a.SearchStr = config.GetString("ldapSearch")

	return nil
}

// Login method is used to check if a couple of user/password is valid in LDAP.
func (a *Authenticator) Login(user, password string) (valid bool, swapuser string, err error) {
	l := &ldap.Conn{}
	l.SetTimeout(time.Second * 10)

	uri := fmt.Sprintf("%s:%d", a.Address, a.Port)

	if a.UseTLS {
		l, err = ldap.DialTLS("tcp", uri, &tls.Config{InsecureSkipVerify: !a.TLSVerify}) // nolint: gosec
	} else {
		l, err = ldap.Dial("tcp", uri)
	}
	if err != nil {
		return false, "", err
	}
	defer l.Close()

	err = l.Bind(a.BindUser, a.BindPassword)
	if err != nil {
		return false, "", err
	}

	searchReq := ldap.NewSearchRequest(
		a.SearchBase, ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		fmt.Sprintf(a.SearchStr, ldap.EscapeFilter(user)),
		[]string{"dn"},
		nil,
	)

	sr, err := l.Search(searchReq)
	if err != nil {
		return false, "", err
	}

	if len(sr.Entries) > 1 {
		return false, "", errors.New("too many user entries returned")
	} else if len(sr.Entries) == 0 {
		return false, "", errors.New("user not found")
	}

	userdn := sr.Entries[0].DN

	// Bind as the user to verify their password
	err = l.Bind(userdn, password)
	if err != nil {
		return false, "", err
	}

	return true, "", nil
}
