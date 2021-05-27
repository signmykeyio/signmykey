package oidcropc

import (
	"context"
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
	OIDCTokenEndpoint string
	OIDCClientID      string
	OIDCClientSecret  string
}

type oidcLogin struct {
	User     string `json:"user" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type oidcTokenResponse struct {
	Token string `json:"access_token"`
	Error string `json:"error"`
}

// OIDCToken represents a OIDC userinfo token
type OIDCToken string

// OIDCTokenKeyType represents an OIDC context key type
type OIDCTokenKeyType string

// OIDCTokenKey represents an OIDC context key
const OIDCTokenKey OIDCTokenKeyType = "oidcToken"

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

	return nil
}

// Login method is used to check if a couple of user/password is valid in OIDC.
func (a *Authenticator) Login(ctx context.Context, payload []byte) (resultCtx context.Context, valid bool, id string, err error) {

	var login oidcLogin
	err = json.Unmarshal(payload, &login)
	if err != nil {
		log.Errorf("json unmarshaling failed: %s", err)
		return ctx, false, "", fmt.Errorf("JSON unmarshaling failed: %w", err)
	}

	v := url.Values{}
	v.Set("grant_type", "password")
	v.Add("username", login.User)
	v.Add("password", login.Password)
	v.Add("client_id", a.OIDCClientID)
	v.Add("client_secret", a.OIDCClientSecret)

	oidcPayload := strings.NewReader(v.Encode())

	reqToken, err := http.NewRequest("POST", a.OIDCTokenEndpoint, oidcPayload)
	if err != nil {
		return ctx, false, "", err
	}

	reqToken.Header.Add("content-type", "application/x-www-form-urlencoded")

	client := http.Client{Timeout: time.Second * 10}
	resToken, err := client.Do(reqToken)
	if err != nil {
		return ctx, false, "", err
	}
	defer resToken.Body.Close()

	bodyToken, err := ioutil.ReadAll(resToken.Body)
	if err != nil {
		return ctx, false, "", errors.New("can't read body")
	}

	tokenRes := oidcTokenResponse{}
	err = json.Unmarshal(bodyToken, &tokenRes)

	if err != nil {
		return ctx, false, "", err
	}

	// return OAuth 2.0 error response if any
	if resToken.StatusCode != 200 {
		return ctx, false, "", errors.New(tokenRes.Error)
	}

	// exit if OIDC Token is empty
	if tokenRes.Token == "" {
		return ctx, false, "", nil
	}

	return context.WithValue(ctx, OIDCTokenKey, OIDCToken(tokenRes.Token)), true, fmt.Sprintf("oidc-%s", login.User), nil
}
