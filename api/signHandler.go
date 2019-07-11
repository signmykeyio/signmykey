package api

import (
	"io/ioutil"
	"net/http"

	"github.com/go-chi/render"
	log "github.com/sirupsen/logrus"
)

func signHandler(w http.ResponseWriter, r *http.Request) {

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Errorf("failed to read body: %s", err)
		render.Status(r, 400)
		render.JSON(w, r, map[string]string{"error": "failed to read body"})
		return
	}

	ctx, valid, id, err := config.Auth.Login(r.Context(), body)
	if !valid {
		log.Errorf("login failed: %s", err)
		render.Status(r, 401)
		render.JSON(w, r, map[string]string{"error": "login failed"})
		return
	}

	ctx, principals, err := config.Princs.Get(ctx, body)
	if err != nil {
		log.Errorf("error getting list of principals: %s", err)
		render.Status(r, 401)
		render.JSON(w, r, map[string]string{"error": "error getting list of principals"})
		return
	}

	cert, err := config.Signer.Sign(ctx, body, id, principals)
	if err != nil {
		log.Errorf("server error during key signing: %s", err)
		render.Status(r, 400)
		render.JSON(w, r, map[string]string{"error": "unknown server error during key signing"})
		return
	}

	log.Debugf("certificate: %s", cert)
	render.JSON(w, r, map[string]string{"certificate": cert})
}
