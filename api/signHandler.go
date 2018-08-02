package api

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/go-chi/render"
	"github.com/pablo-ruth/signmykey/builtin/signer"
	log "github.com/sirupsen/logrus"
)

func signHandler(w http.ResponseWriter, r *http.Request) {
	var login Login

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Errorf("failed to read body: %s", err)
		render.Status(r, 400)
		render.JSON(w, r, map[string]string{"error": "failed to read body"})
		return
	}

	err = json.Unmarshal(body, &login)
	if err != nil {
		log.Errorf("json unmarshaling failed: %s", err)
		render.Status(r, 400)
		render.JSON(w, r, map[string]string{"error": "JSON unmarshaling failed"})
		return
	}

	err = login.Validate()
	if err != nil {
		log.Errorf("missing field(s) in signing request: %s", err)
		render.Status(r, 400)
		render.JSON(w, r, map[string]string{"error": "missing field(s) in signing request"})
		return
	}

	valid, err := config.Auth.Login(login.User, login.Password)
	if !valid {
		log.Errorf("login failed: %s", err)
		render.Status(r, 401)
		render.JSON(w, r, map[string]string{"error": "login failed"})
		return
	}

	principals, err := config.Princs.Get(login.User)
	if err != nil {
		log.Errorf("error getting list of principals: %s", err)
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
		log.Errorf("server error during key signing: %s", err)
		render.Status(r, 400)
		render.JSON(w, r, map[string]string{"error": "unknown server error during key signing"})
		return
	}

	log.Debugf("certificate: %s", cert)
	render.JSON(w, r, map[string]string{"certificate": cert})
}
