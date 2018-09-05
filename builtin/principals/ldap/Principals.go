package ldap

import (
	"crypto/tls"
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/spf13/viper"
	ldap "gopkg.in/ldap.v2"
)

// Principals struct represents LDAP options for getting principals list from LDAP.
type Principals struct {
	Address         string
	Port            int
	BindUser        string
	BindPassword    string
	UserSearchBase  string
	UserSearchStr   string
	GroupSearchBase string
	GroupSearchStr  string
	UseTLS          bool
	TLSVerify       bool
	Prefix          string
	TransformCase   string
}

// Init method is used to ingest config of Principals
func (p *Principals) Init(config *viper.Viper) error {
	neededEntries := []string{
		"ldapAddr",
		"ldapPort",
		"ldapTLS",
		"ldapTLSVerify",
		"ldapBindUser",
		"ldapBindPassword",
		"ldapUserBase",
		"ldapUserSearch",
		"ldapGroupBase",
		"ldapGroupSearch",
	}

	var missingEntriesLst []string
	for _, entry := range neededEntries {
		if !config.IsSet(entry) {
			missingEntriesLst = append(missingEntriesLst, entry)
			continue
		}
	}
	if len(missingEntriesLst) > 0 {
		missingEntries := strings.Join(missingEntriesLst, ", ")
		return fmt.Errorf("Missing config entries (%s) for Principals", missingEntries)
	}

	config.SetDefault("transformCase", "none")
	tc := config.GetString("transformCase")
	if tc != "none" && tc != "lower" && tc != "upper" {
		return errors.New("transformCase config entry for Principals must be none, lower or upper")
	}

	p.Address = config.GetString("ldapAddr")
	p.Port = config.GetInt("ldapPort")
	p.UseTLS = config.GetBool("ldapTLS")
	p.TLSVerify = config.GetBool("ldapTLSVerify")
	p.BindUser = config.GetString("ldapBindUser")
	p.BindPassword = config.GetString("ldapBindPassword")
	p.UserSearchBase = config.GetString("ldapUserBase")
	p.UserSearchStr = config.GetString("ldapUserSearch")
	p.GroupSearchBase = config.GetString("ldapGroupBase")
	p.GroupSearchStr = config.GetString("ldapGroupSearch")
	p.Prefix = config.GetString("ldapGroupPrefix")
	p.TransformCase = tc

	return nil
}

// Get method is used to get the list of principals associated to a specific user.
func (p Principals) Get(user string) (principals []string, err error) {
	l, err := getLDAPConn(p)
	if err != nil {
		return principals, err
	}
	defer l.Close()

	userSearchReq := ldap.NewSearchRequest(
		p.UserSearchBase, ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		fmt.Sprintf(p.UserSearchStr, user),
		[]string{},
		nil,
	)

	usr, err := l.Search(userSearchReq)
	if err != nil {
		return principals, err
	}

	if len(usr.Entries) > 1 {
		return principals, errors.New("too many user entries returned")
	} else if len(usr.Entries) == 0 {
		return principals, errors.New("user not found")
	}

	groupSearchRequest := ldap.NewSearchRequest(
		p.GroupSearchBase,
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		fmt.Sprintf(p.GroupSearchStr, usr.Entries[0].DN),
		[]string{},
		nil,
	)

	gsr, err := l.Search(groupSearchRequest)
	if err != nil {
		return principals, err
	}

	if len(gsr.Entries) == 0 {
		return principals, errors.New("no group found for this user")
	}

	for _, group := range gsr.Entries {
		principals = append(principals, group.DN)
	}

	principals = getCN(principals)
	principals = filterByPrefix(p.Prefix, principals)
	principals = transformCase(p.TransformCase, principals)

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

func transformCase(transform string, list []string) []string {
	principals := []string{}

	if transform == "lower" {
		for _, str := range list {
			principals = append(principals, strings.ToLower(str))
		}

		return principals
	}

	if transform == "upper" {
		for _, str := range list {
			principals = append(principals, strings.ToUpper(str))
		}

		return principals
	}

	return list
}

func getLDAPConn(p Principals) (l *ldap.Conn, err error) {
	l = &ldap.Conn{}
	l.SetTimeout(time.Second * 10)

	uri := fmt.Sprintf("%s:%d", p.Address, p.Port)

	if p.UseTLS {
		l, err = ldap.DialTLS("tcp", uri, &tls.Config{InsecureSkipVerify: !p.TLSVerify}) // nolint: gas
	} else {
		l, err = ldap.Dial("tcp", uri)
	}
	if err != nil {
		return l, err
	}

	err = l.Bind(p.BindUser, p.BindPassword)

	return l, err
}
