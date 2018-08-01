package vault

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

// Authenticator struct represents Vault options for SMK Authentication.
type Authenticator struct {
	Address  string
	Port     int
	UseTLS   bool
	Path     string
	fullAddr string
}

// Init method is used to ingest config of Authenticator
func (v *Authenticator) Init(config map[string]string) error {
	neededEntries := []string{
		"vaultAddr",
		"vaultPort",
		"vaultTLS",
		"vaultPath",
	}

	for _, entry := range neededEntries {
		if _, ok := config[entry]; !ok {
			return fmt.Errorf("Config entry %s missing for Authenticator", entry)
		}
	}

	// Conversions
	port, err := strconv.Atoi(config["vaultPort"])
	if err != nil {
		return err
	}
	useTLS, err := strconv.ParseBool(config["vaultTLS"])
	if err != nil {
		return err
	}

	v.Address = config["vaultAddr"]
	v.Port = port
	v.UseTLS = useTLS
	v.Path = config["vaultPath"]

	var scheme string
	if v.UseTLS {
		scheme = "https"
	} else {
		scheme = "http"
	}
	v.fullAddr = fmt.Sprintf("%s://%s:%d/v1", scheme, v.Address, v.Port)

	return nil
}

// Login method is used to check if a couple of user/password is valid in Vault.
func (v Authenticator) Login(user, password string) (valid bool, err error) {
	if len(user) == 0 {
		return false, fmt.Errorf("empty username")
	}
	if len(password) == 0 {
		return false, fmt.Errorf("empty password")
	}

	data, err := json.Marshal(map[string]string{"password": password})
	if err != nil {
		return false, err
	}

	resp, err := http.Post(
		fmt.Sprintf("%s/auth/%s/login/%s", v.fullAddr, v.Path, user),
		"application/json",
		bytes.NewBuffer(data))

	if err != nil {
		return false, err
	}
	defer resp.Body.Close() // nolint: errcheck

	if resp.StatusCode == 400 {
		return false, fmt.Errorf("invalid username or password")
	}
	if resp.StatusCode == 500 {
		return false, fmt.Errorf("Vault internal server error")
	}
	if resp.StatusCode != 200 {
		return false, fmt.Errorf("unknown error during authentication")
	}

	return true, nil
}
