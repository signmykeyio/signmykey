package api

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
	"github.com/signmykeyio/signmykey/builtin/signer"
	"github.com/sirupsen/logrus"
)

func signHandler(w http.ResponseWriter, r *http.Request) {
	var login Login

	log := r.Context().Value("logger").(*logrus.Logger)
	logger := log.WithField("app", "http")
	reqID := middleware.GetReqID(r.Context())

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		logger.Warnf("[%s] failed to read body: %s", reqID, err)
		render.Status(r, 400)
		render.JSON(w, r, map[string]string{"error": "failed to read body"})
		return
	}

	err = json.Unmarshal(body, &login)
	if err != nil {
		logger.Warnf("[%s] json unmarshaling failed: %s", reqID, err)
		render.Status(r, 400)
		render.JSON(w, r, map[string]string{"error": "JSON unmarshaling failed"})
		return
	}

	err = login.Validate()
	if err != nil {
		logger.Warnf("[%s] missing field(s) in signing request: %s", reqID, err)
		render.Status(r, 400)
		render.JSON(w, r, map[string]string{"error": "missing field(s) in signing request"})
		return
	}

	valid, err := config.Auth.Login(login.User, login.Password)
	if !valid {
		logger.Warnf("[%s] login failed: %s", reqID, err)
		render.Status(r, 401)
		render.JSON(w, r, map[string]string{"error": "login failed"})
		return
	}

	principals, err := config.Princs.Get(login.User)
	if err != nil {
		logger.Warnf("[%s] error getting list of principals: %s", reqID, err)
		render.Status(r, 401)
		render.JSON(w, r, map[string]string{"error": "error getting list of principals"})
		return
	}

	req := signer.CertReq{
		Key:        login.PubKey,
		ID:         login.User,
		Principals: principals,
	}
	cert, err := config.Signer.Sign(req)
	if err != nil {
		logger.Warnf("[%s] server error during key signing: %s", reqID, err)
		render.Status(r, 400)
		render.JSON(w, r, map[string]string{"error": "unknown server error during key signing"})
		return
	}

	render.JSON(w, r, map[string]string{"certificate": cert})
}
