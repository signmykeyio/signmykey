package api

import (
	"errors"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
	"github.com/signmykeyio/signmykey/client"
	"github.com/sirupsen/logrus"

	builtinPrincs "github.com/signmykeyio/signmykey/builtin/principals"
)

func signHandler(w http.ResponseWriter, r *http.Request) {

	log := r.Context().Value(RequestLoggerKey).(*logrus.Logger)
	reqID := middleware.GetReqID(r.Context())

	logger := log.WithFields(logrus.Fields{
		"ctx":     "api",
		"handler": "sign",
		"req_id":  reqID,
	})

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		logger.WithError(err).Error("Reading signing request body")
		render.Status(r, 400)
		render.JSON(w, r, map[string]string{"error": "failed to read body"})
		return
	}

	ctx, valid, id, err := config.Auth.Login(r.Context(), body)
	if !valid {
		logger.WithError(err).Error("Authenticating user")
		render.Status(r, 401)
		render.JSON(w, r, map[string]string{"error": "login failed"})
		return
	}
	logger = logger.WithField("user", id)
	logger.Info("User authenticated")

	principals := []string{}
	for _, princsProvider := range config.Princs {
		_, princs, err := princsProvider.Get(ctx, body)
		if err != nil {
			// this is not critical error, next provider cat return principals
			var principalsNotFoundError *builtinPrincs.NotFoundError
			if errors.As(err, &principalsNotFoundError) {
				logger.Info(err.Error())
				continue
			}
			logger.WithError(err).Error("Getting list of user principals")
			render.Status(r, 401)
			render.JSON(w, r, map[string]string{"error": "error getting list of principals"})
			return
		}

		principals = append(principals, princs...)
	}

	if len(principals) == 0 {
		logger.Error("No principals found")
		render.Status(r, 401)
		render.JSON(w, r, map[string]string{"error": "no principals found"})
		return
	}

	logger = logger.WithField("principals", principals)
	logger.Info("User principals retrieved")

	cert, err := config.Signer.Sign(ctx, body, id, principals)
	if err != nil {
		logger.WithError(err).Error("Generating SSH certificate")
		render.Status(r, 400)
		render.JSON(w, r, map[string]string{"error": "unknown server error during key signing"})
		return
	}

	_, before, _ := client.CertInfo(cert)
	logger.WithField("expire", time.Unix(int64(before), 0)).Info("SSH certificate generated")

	render.JSON(w, r, map[string]string{"certificate": cert})
}
