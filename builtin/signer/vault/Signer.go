package vault

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/signmykeyio/signmykey/builtin/signer"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// Signer struct represents Hashicorp Vault options for signing SSH Key.
type Signer struct {
	Address  string
	Port     int
	UseTLS   bool
	RoleID   string
	SecretID string
	Path     string
	Role     string
	SignTTL  string
	fullAddr string
}

type vaultSignReq struct {
	PubKey string `json:"public_key" binding:"required"`
}

// Init method is used to ingest config of Signer
func (v *Signer) Init(config *viper.Viper) error {
	neededEntries := []string{
		"vaultAddr",
		"vaultPort",
		"vaultTLS",
		"vaultRoleID",
		"vaultSecretID",
		"vaultPath",
		"vaultRole",
		"vaultSignTTL",
	}

	for _, entry := range neededEntries {
		if !config.IsSet(entry) {
			return fmt.Errorf("Config entry %s missing for Signer", entry)
		}
	}

	v.Address = config.GetString("vaultAddr")
	v.Port = config.GetInt("vaultPort")
	v.UseTLS = config.GetBool("vaultTLS")
	v.RoleID = config.GetString("vaultRoleID")
	v.SecretID = config.GetString("vaultSecretID")
	v.Path = config.GetString("vaultPath")
	v.Role = config.GetString("vaultRole")
	v.SignTTL = config.GetString("vaultSignTTL")

	var scheme string
	if v.UseTLS {
		scheme = "https"
	} else {
		scheme = "http"
	}
	v.fullAddr = fmt.Sprintf("%s://%s:%d/v1", scheme, v.Address, v.Port)

	return nil
}

// ReadCA method read CA public cert from Hashicorp Vault backend
func (v Signer) ReadCA() (string, error) {
	token, err := v.getToken()
	if err != nil {
		return "", fmt.Errorf("error getting auth token: %w", err)
	}

	client := http.Client{Timeout: time.Second * 10}
	req, err := http.NewRequest(
		"GET",
		fmt.Sprintf("%s/%s/config/ca", v.fullAddr, v.Path),
		nil,
	)
	if err != nil {
		return "", fmt.Errorf("error creating new httprequest: %w", err)
	}
	req.Header.Add("X-Vault-Token", token)

	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed during Get request: %w", err)
	}
	defer resp.Body.Close() // nolint: errcheck

	if resp.StatusCode != 200 {
		return "", fmt.Errorf("unknown error from Vault with status code %d", resp.StatusCode)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error reading response body: %w", err)
	}

	var caResp map[string]interface{}
	err = json.Unmarshal(body, &caResp)
	if err != nil {
		return "", fmt.Errorf("error unmarshaling vault response: %w", err)
	}
	val, ok := caResp["data"].(map[string]interface{})
	if !ok {
		return "", fmt.Errorf("CA public key not found in Vault response")
	}
	publicKey, ok := val["public_key"].(string)
	if !ok {
		return "", fmt.Errorf("CA public key not found in Vault response")
	}

	// Remove newline
	publicKey = strings.Replace(publicKey, "\n", "", 1)

	return publicKey, nil
}

// Sign method is used to sign passed SSH Key.
func (v Signer) Sign(ctx context.Context, payload []byte, id string, principals []string) (cert string, err error) {

	var signReq vaultSignReq
	err = json.Unmarshal(payload, &signReq)
	if err != nil {
		log.Errorf("json unmarshaling failed: %s", err)
		return "", fmt.Errorf("JSON unmarshaling failed: %w", err)
	}

	token, err := v.getToken()
	if err != nil {
		return "", fmt.Errorf("error getting auth token: %w", err)
	}

	certreq := signer.CertReq{
		Key:        signReq.PubKey,
		ID:         id,
		Principals: principals,
	}
	data, err := json.Marshal(map[string]string{
		"key_id":           certreq.ID,
		"public_key":       certreq.Key,
		"valid_principals": strings.Join(certreq.Principals, ","),
		"ttl":              v.SignTTL,
	})
	if err != nil {
		return "", fmt.Errorf("marshaling of sign request payload failed: %w", err)
	}

	client := http.Client{Timeout: time.Second * 10}
	req, err := http.NewRequest(
		"POST",
		fmt.Sprintf("%s/%s/sign/%s", v.fullAddr, v.Path, v.Role),
		bytes.NewBuffer(data),
	)
	if err != nil {
		return "", fmt.Errorf("error creating new httprequest: %w", err)
	}
	req.Header.Add("X-Vault-Token", token)

	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed during POST request: %w", err)
	}
	defer resp.Body.Close() // nolint: errcheck

	if resp.StatusCode != 200 {
		return "", fmt.Errorf("unknown error from Vault with status code %d", resp.StatusCode)
	}

	signedKey, err := extractSignedKey(resp)

	return signedKey, err
}

func extractSignedKey(resp *http.Response) (string, error) {
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error reading response body: %w", err)
	}

	var signResp map[string]interface{}
	err = json.Unmarshal(body, &signResp)
	if err != nil {
		return "", fmt.Errorf("error unmarshaling Vault response: %w", err)
	}
	val, ok := signResp["data"].(map[string]interface{})
	if !ok {
		return "", fmt.Errorf("Signed key not found in Vault response")
	}
	signedKey, ok := val["signed_key"].(string)
	if !ok {
		return "", fmt.Errorf("Signed key not found in Vault response")
	}

	return signedKey, nil
}

func (v Signer) getToken() (string, error) {

	data, err := json.Marshal(map[string]string{
		"role_id":   v.RoleID,
		"secret_id": v.SecretID,
	})
	if err != nil {
		return "", err
	}

	client := http.Client{Timeout: time.Second * 10}
	resp, err := client.Post(
		fmt.Sprintf("%s/auth/approle/login", v.fullAddr),
		"application/json",
		bytes.NewBuffer(data))

	if err != nil {
		return "", err
	}
	defer resp.Body.Close() // nolint: errcheck

	if resp.StatusCode == 400 {
		return "", fmt.Errorf("invalid role_id or secret_id")
	}
	if resp.StatusCode == 500 {
		return "", fmt.Errorf("Vault internal server error")
	}
	if resp.StatusCode != 200 {
		return "", fmt.Errorf("unknown error during authentication with status code %d", resp.StatusCode)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var loginResp map[string]interface{}
	err = json.Unmarshal(body, &loginResp)
	if err != nil {
		return "", err
	}
	val, ok := loginResp["auth"].(map[string]interface{})
	if !ok {
		return "", fmt.Errorf("client token not found in AppRole login response")
	}
	token, ok := val["client_token"].(string)
	if !ok {
		return "", fmt.Errorf("client token not found in AppRole login response")
	}

	return token, nil
}
