package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestSignHandler(t *testing.T) {
	type JSONResponse map[string]interface{}

	cases := []struct {
		method      string
		url         string
		code        int
		payload     []byte
		response    interface{}
		contentType string
	}{
		{"GET", "/v1/sign", 405, []byte(""), JSONResponse(nil), ""},
		{"PUT", "/v1/sign", 405, []byte(""), JSONResponse(nil), ""},
		{"PATCH", "/v1/sign", 405, []byte(""), JSONResponse(nil), ""},
		{"DELETE", "/v1/sign", 405, []byte(""), JSONResponse(nil), ""},
		{
			"POST", "/v1/sign", 400,
			[]byte("{\"user\":\"test\"}"),
			JSONResponse{"error": "missing field(s) in signing request"},
			"application/json",
		},
		{
			"POST", "/v1/sign", 401,
			[]byte("{\"user\":\"baduser\",\"password\":\"badpassword\",\"public_key\":\"goodkey\"}"),
			JSONResponse{"error": "login failed"},
			"application/json",
		},
		{
			"POST", "/v1/sign", 401,
			[]byte("{\"user\":\"testuser\",\"password\":\"badpassword\",\"public_key\":\"goodkey\"}"),
			JSONResponse{"error": "login failed"},
			"application/json",
		},
		{
			"POST", "/v1/sign", 401,
			[]byte("{\"user\":\"emptyprincsuser\",\"password\":\"testpassword\",\"public_key\":\"goodkey\"}"),
			JSONResponse{"error": "error getting list of principals"},
			"application/json",
		},
		{
			"POST", "/v1/sign", 400,
			[]byte("{\"user\":\"testuser\",\"password\":\"testpassword\",\"public_key\":\"badkey\"}"),
			JSONResponse{"error": "unknown server error during key signing"},
			"application/json",
		},
		{
			"POST", "/v1/sign", 200,
			[]byte("{\"user\":\"testuser\",\"password\":\"testpassword\",\"public_key\":\"goodkey\"}"),
			JSONResponse{"certificate": "goodcert"},
			"application/json",
		},
	}

	config = Config{
		Auth:   &authMock{},
		Princs: &princsMock{},
		Signer: &signerMock{},
	}
	router := Router()

	for _, c := range cases {
		w := httptest.NewRecorder()
		mj, _ := json.Marshal(c.payload)
		req, _ := http.NewRequest(c.method, c.url, bytes.NewBuffer(mj))
		router.ServeHTTP(w, req)

		var response JSONResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		if err != nil {
			assert.Fail(t, err.Error())
		}

		assert.Equal(t, w.Code, c.code)
		assert.Equal(t, c.response, response)
		assert.Contains(t, w.Header().Get("Content-Type"), c.contentType)
	}
}

type authMock struct{}

func (a authMock) Login(ctx context.Context, payload []byte) (context.Context, bool, string, error) {
	var login struct {
		User     string `json:"user" binding:"required"`
		Password string `json:"password" binding:"required"`
	}
	err := json.Unmarshal(payload, &login)
	if err != nil {
		log.Errorf("json unmarshaling failed: %s", err)
		return ctx, false, "", fmt.Errorf("JSON unmarshaling failed: %s", err)
	}

	if login.User != "testuser" && login.User != "emptyprincsuser" {
		return ctx, false, "", fmt.Errorf("unknown username")
	}

	if login.Password != "testpassword" {
		return ctx, false, "", fmt.Errorf("invalid password")
	}

	return ctx, true, "", nil
}

func (a authMock) Init(config *viper.Viper) error {
	return nil
}

type princsMock struct{}

func (p princsMock) Init(config *viper.Viper) error {
	return nil
}

func (p princsMock) Get(ctx context.Context, payload []byte) (context.Context, []string, error) {
	var login struct {
		User     string `json:"user" binding:"required"`
		Password string `json:"password" binding:"required"`
	}
	err := json.Unmarshal(payload, &login)
	if err != nil {
		log.Errorf("json unmarshaling failed: %s", err)
		return ctx, []string{}, fmt.Errorf("JSON unmarshaling failed: %s", err)
	}

	if login.User == "emptyprincsuser" {
		return ctx, []string{}, fmt.Errorf("empty list of principals")
	}

	return ctx, []string{"root", "user"}, nil
}

type signerMock struct{}

func (s signerMock) Init(config *viper.Viper) error {
	return nil
}

func (s signerMock) ReadCA() (string, error) {
	return "", nil
}

func (s signerMock) Sign(ctx context.Context, payload []byte, id string, principals []string) (string, error) {
	var pubkey struct {
		PubKey string `json:"public_key" binding:"required"`
	}
	err := json.Unmarshal(payload, &pubkey)
	if err != nil {
		log.Errorf("json unmarshaling failed: %s", err)
		return "", fmt.Errorf("JSON unmarshaling failed: %s", err)
	}

	if pubkey.PubKey == "goodkey" {
		return "goodcert", nil
	}

	if pubkey.PubKey == "badkey" {
		return "", fmt.Errorf("bad key format")
	}

	return "", fmt.Errorf("failed to sign key")
}
