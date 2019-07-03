package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/signmykeyio/signmykey/builtin/signer"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestSignHandler(t *testing.T) {
	type JSONResponse map[string]interface{}

	cases := []struct {
		method      string
		url         string
		code        int
		payload     Login
		response    interface{}
		contentType string
	}{
		{"GET", "/v1/sign", 405, Login{}, JSONResponse(nil), ""},
		{"PUT", "/v1/sign", 405, Login{}, JSONResponse(nil), ""},
		{"PATCH", "/v1/sign", 405, Login{}, JSONResponse(nil), ""},
		{"DELETE", "/v1/sign", 405, Login{}, JSONResponse(nil), ""},
		{
			"POST", "/v1/sign", 400,
			Login{User: "test"},
			JSONResponse{"error": "missing field(s) in signing request"},
			"application/json",
		},
		{
			"POST", "/v1/sign", 401,
			Login{User: "baduser", Password: "badpassword", PubKey: "goodkey"},
			JSONResponse{"error": "login failed"},
			"application/json",
		},
		{
			"POST", "/v1/sign", 401,
			Login{User: "testuser", Password: "badpassword", PubKey: "goodkey"},
			JSONResponse{"error": "login failed"},
			"application/json",
		},
		{
			"POST", "/v1/sign", 401,
			Login{User: "emptyprincsuser", Password: "testpassword", PubKey: "goodkey"},
			JSONResponse{"error": "error getting list of principals"},
			"application/json",
		},
		{
			"POST", "/v1/sign", 400,
			Login{User: "testuser", Password: "testpassword", PubKey: "badkey"},
			JSONResponse{"error": "unknown server error during key signing"},
			"application/json",
		},
		{
			"POST", "/v1/sign", 200,
			Login{User: "testuser", Password: "testpassword", PubKey: "goodkey"},
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
		json.Unmarshal(w.Body.Bytes(), &response)

		assert.Equal(t, w.Code, c.code)
		assert.Equal(t, c.response, response)
		assert.Contains(t, w.Header().Get("Content-Type"), c.contentType)
	}
}

type authMock struct{}

func (a authMock) Login(user, password string) (bool, string, error) {
	if user != "testuser" && user != "emptyprincsuser" {
		return false, "", fmt.Errorf("unknown username")
	}

	if password != "testpassword" {
		return false, "", fmt.Errorf("invalid password")
	}

	return true, "", nil
}

func (a authMock) Init(config *viper.Viper) error {
	return nil
}

type princsMock struct{}

func (p princsMock) Init(config *viper.Viper) error {
	return nil
}

func (p princsMock) Get(user string) ([]string, error) {
	if user == "emptyprincsuser" {
		return []string{}, fmt.Errorf("empty list of principals")
	}

	return []string{"root", "user"}, nil
}

type signerMock struct{}

func (s signerMock) Init(config *viper.Viper) error {
	return nil
}

func (s signerMock) ReadCA() (string, error) {
	return "", nil
}

func (s signerMock) Sign(req signer.CertReq) (string, error) {
	if req.Key == "goodkey" {
		return "goodcert", nil
	}

	if req.Key == "badkey" {
		return "", fmt.Errorf("bad key format")
	}

	return "", fmt.Errorf("failed to sign key")
}
