package ldap

import (
	"crypto/tls"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"gopkg.in/ldap.v2"
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
func (a *Authenticator) Init(config map[string]string) error {
	neededEntries := []string{
		"ldap-addr",
		"ldap-port",
		"ldap-tls",
		"ldap-tls-verify",
		"ldap-bind-user",
		"ldap-bind-password",
		"ldap-base",
		"ldap-search",
	}

	var missingEntriesLst []string
	for _, entry := range neededEntries {
		if _, ok := config[entry]; !ok {
			missingEntriesLst = append(missingEntriesLst, entry)
			continue
		}

		if len(config[entry]) == 0 {
			missingEntriesLst = append(missingEntriesLst, entry)
		}
	}
	if len(missingEntriesLst) > 0 {
		missingEntries := strings.Join(missingEntriesLst, ", ")
		return fmt.Errorf("Missing config entries (%s) for Authenticator", missingEntries)
	}

	// Conversions
	port, err := strconv.Atoi(config["ldap-port"])
	if err != nil {
		return err
	}
	useTLS, err := strconv.ParseBool(config["ldap-tls"])
	if err != nil {
		return err
	}
	tlsVerify, err := strconv.ParseBool(config["ldap-tls-verify"])
	if err != nil {
		return err
	}

	a.Address = config["ldap-addr"]
	a.Port = port
	a.UseTLS = useTLS
	a.BindUser = config["ldap-bind-user"]
	a.BindPassword = config["ldap-bind-password"]
	a.SearchBase = config["ldap-base"]
	a.SearchStr = config["ldap-search"]
	a.TLSVerify = tlsVerify

	return nil
}

// Login method is used to check if a couple of user/password is valid in LDAP.
func (a *Authenticator) Login(user, password string) (valid bool, err error) {
	var l *ldap.Conn
	uri := fmt.Sprintf("%s:%d", a.Address, a.Port)

	if a.UseTLS {
		l, err = ldap.DialTLS("tcp", uri, &tls.Config{InsecureSkipVerify: !a.TLSVerify}) // nolint: gas
	} else {
		l, err = ldap.Dial("tcp", uri)
	}
	if err != nil {
		return false, err
	}
	defer l.Close()

	err = l.Bind(a.BindUser, a.BindPassword)
	if err != nil {
		return false, err
	}

	searchReq := ldap.NewSearchRequest(
		a.SearchBase, ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		fmt.Sprintf(a.SearchStr, user),
		[]string{"dn"},
		nil,
	)

	sr, err := l.Search(searchReq)
	if err != nil {
		return false, err
	}

	if len(sr.Entries) > 1 {
		return false, errors.New("too many user entries returned")
	} else if len(sr.Entries) == 0 {
		return false, errors.New("user not found")
	}

	userdn := sr.Entries[0].DN

	// Bind as the user to verify their password
	err = l.Bind(userdn, password)
	if err != nil {
		return false, err
	}

	return true, nil
}
