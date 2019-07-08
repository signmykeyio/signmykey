package oidcropc

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
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

type oidcToken struct {
	Token string `json:"access_token"`
	Error string `json:"error"`
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

// Login method is used to check if a couple of user/password is valid in OIDC.
func (a *Authenticator) Login(user, password string) (valid bool, swapuser string, err error) {

	v := url.Values{}
	v.Set("grant_type", "password")
	v.Add("username", user)
	v.Add("password", password)
	v.Add("client_id", a.OIDCClientID)
	v.Add("client_secret", a.OIDCClientSecret)

	payload := strings.NewReader(v.Encode())

	reqToken, err := http.NewRequest("POST", a.OIDCTokenEndpoint, payload)
	if err != nil {
		return false, "", err
	}

	reqToken.Header.Add("content-type", "application/x-www-form-urlencoded")

	client := http.Client{Timeout: time.Second * 10}
	resToken, err := client.Do(reqToken)
	if err != nil {
		return false, "", err
	}

	defer resToken.Body.Close()

	bodyToken, err := ioutil.ReadAll(resToken.Body)
	if err != nil {
		return false, "", errors.New("can't read body")
	}

	oidcToken1 := oidcToken{}
	err = json.Unmarshal(bodyToken, &oidcToken1)

	log.Debugf("OIDC Token: %s", oidcToken1.Token)

	if err != nil {
		return false, "", err
	}

	// return OAuth 2.0 error response if any
	if resToken.StatusCode == 400 {
		return false, "", errors.New(oidcToken1.Error)
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
