package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPingHandler(t *testing.T) {
	type JSONResponse map[string]interface{}

	cases := []struct {
		method      string
		url         string
		code        int
		response    interface{}
		contentType string
	}{
		{"GET", "/v1/ping", 200, JSONResponse{"message": "pong"}, "application/json"},
		{"POST", "/v1/ping", 405, JSONResponse(nil), ""},
		{"PUT", "/v1/ping", 405, JSONResponse(nil), ""},
		{"PATCH", "/v1/ping", 405, JSONResponse(nil), ""},
		{"DELETE", "/v1/ping", 405, JSONResponse(nil), ""},
	}

	router := Router()

	for _, c := range cases {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(c.method, c.url, nil)
		router.ServeHTTP(w, req)

		if c.code != w.Code {
			assert.FailNowf(t, "bad status code", "expected %d, have %d", c.code, w.Code)
		}

		if c.code != 200 {
			continue
		}

		var response JSONResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		if err != nil {
			assert.Failf(t, "Error not expected", err.Error())
		}

		assert.Equal(t, c.response, response)
		assert.Contains(t, w.Header().Get("Content-Type"), c.contentType)
	}
}
