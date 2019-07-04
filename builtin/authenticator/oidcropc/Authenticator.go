package oidcropc

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// Authenticator struct represents OIDC options for SMK Authentication.
type Authenticator struct {
	OIDCTokenEndpoint      string
	OIDCClientID           string
	OIDCClientSecret       string
	OIDCAlternatePrincipal bool
}

// Init method is used to ingest config of Authenticator
func (a *Authenticator) Init(config *viper.Viper) error {
	neededEntries := []string{
		"oidcTokenEndpoint",
		"oidcClientID",
		"oidcClientSecret",
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

	a.OIDCTokenEndpoint = config.GetString("oidcTokenEndpoint")
	a.OIDCClientID = config.GetString("oidcClientID")
	a.OIDCClientSecret = config.GetString("oidcClientSecret")
	a.OIDCAlternatePrincipal = config.GetBool("oidcAlternatePrincipal")

	return nil
}

type oidcToken struct {
	Token string `json:"access_token"`
}

// Login method is used to check if a couple of user/password is valid in LDAP.
func (a *Authenticator) Login(user, password string) (valid bool, swapuser string, err error) {

	payload := strings.NewReader("grant_type=password&" +
		"username=" +
		user +
		"&password=" +
		password +
		"&client_id=" +
		a.OIDCClientID +
		"&client_secret=" +
		a.OIDCClientSecret)

	reqToken, _ := http.NewRequest("POST", a.OIDCTokenEndpoint, payload)

	reqToken.Header.Add("content-type", "application/x-www-form-urlencoded")

	client := http.Client{Timeout: time.Second * 10}
	resToken, errResToken := client.Do(reqToken)
	if errResToken != nil {
		return false, "", errors.New("bad password")
	}

	defer resToken.Body.Close()

	bodyToken, _ := ioutil.ReadAll(resToken.Body)

	oidcToken1 := oidcToken{}
	jsonTokenErr := json.Unmarshal(bodyToken, &oidcToken1)

	log.Debugf("OIDC Token: %s", oidcToken1.Token)

	if jsonTokenErr != nil {
		return false, "", jsonTokenErr
	}

	// exit if OIDC Token is empty
	if oidcToken1.Token == "" {
		return false, "", nil
	}

	// if not using OIDC Principals flow
	if a.OIDCAlternatePrincipal {
		swapuser = ""
	} else {
		swapuser = oidcToken1.Token
	}

	return true, swapuser, nil
}
