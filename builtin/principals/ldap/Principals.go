package ldap

import (
	"crypto/tls"
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"gopkg.in/ldap.v2"
)

// Principals struct represents LDAP options for getting principals list from LDAP.
type Principals struct {
	Address      string
	Port         int
	BindUser     string
	BindPassword string
	SearchBase   string
	SearchStr    string
	UseTLS       bool
	TLSVerify    bool
	Prefix       string
}

// Init method is used to ingest config of Principals
func (p *Principals) Init(config map[string]string) error {
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
		return fmt.Errorf("Missing config entries (%s) for Principals", missingEntries)
	}

	// Conversions
	port, err := strconv.Atoi(config["ldapPort"])
	if err != nil {
		return err
	}
	useTLS, err := strconv.ParseBool(config["ldapTLS"])
	if err != nil {
		return err
	}
	tlsVerify, err := strconv.ParseBool(config["ldapTLSVerify"])
	if err != nil {
		return err
	}

	p.Address = config["ldapAddr"]
	p.Port = port
	p.UseTLS = useTLS
	p.BindUser = config["ldapBindUser"]
	p.BindPassword = config["ldapBindPassword"]
	p.SearchBase = config["ldapBase"]
	p.SearchStr = config["ldapSearch"]
	p.TLSVerify = tlsVerify

	return nil
}

// Get method is used to get the list of principals associated to a specific user.
func (p Principals) Get(user string) (principals []string, err error) {
	var l *ldap.Conn
	uri := fmt.Sprintf("%s:%d", p.Address, p.Port)

	if p.UseTLS {
		l, err = ldap.DialTLS("tcp", uri, &tls.Config{InsecureSkipVerify: !p.TLSVerify}) // nolint: gas
	} else {
		l, err = ldap.Dial("tcp", uri)
	}
	if err != nil {
		return principals, err
	}
	defer l.Close()

	err = l.Bind(p.BindUser, p.BindPassword)
	if err != nil {
		return principals, err
	}

	searchReq := ldap.NewSearchRequest(
		p.SearchBase, ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		fmt.Sprintf(p.SearchStr, user),
		[]string{"memberOf"},
		nil,
	)

	sr, err := l.Search(searchReq)
	if err != nil {
		return principals, err
	}

	if len(sr.Entries) > 1 {
		return principals, errors.New("too many user entries returned")
	} else if len(sr.Entries) == 0 {
		return principals, errors.New("user not found")
	}

	principals = sr.Entries[0].GetAttributeValues("memberOf")
	principals = getCN(principals)
	principals = filterByPrefix(p.Prefix, principals)

	return principals, nil
}

func getCN(list []string) []string {
	var groupRegex = regexp.MustCompile(`^[cC][nN]=(.+?),.*$`)
	cnList := []string{}

	for _, str := range list {
		match := groupRegex.FindStringSubmatch(str)
		if len(match) > 0 && match[1][:1] != "," {
			cnList = append(cnList, match[1])
		}
	}

	return cnList
}

func filterByPrefix(prefix string, list []string) []string {
	principals := []string{}

	for _, str := range list {
		if strings.HasPrefix(str, prefix) {
			principals = append(principals, str[len(prefix):])
		}
	}
	return principals
}
