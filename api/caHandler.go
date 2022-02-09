package api

import (
	"net/http"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/sirupsen/logrus"
)

func caHandler(w http.ResponseWriter, r *http.Request) {

	log := r.Context().Value(RequestLoggerKey).(*logrus.Logger)
	reqID := middleware.GetReqID(r.Context())

	logger := log.WithFields(logrus.Fields{
		"ctx":     "api",
		"handler": "ca",
		"req_id":  reqID,
	})

	publicKey, err := config.Signer.ReadCA(r.Context())
	if err != nil {
		logger.WithError(err).Error("Getting SSH CA certificate")
		render.Status(r, 500)
		render.JSON(w, r, map[string]string{"errors": "error getting CA certificate"})
		return
	}

	render.JSON(w, r, map[string]string{"public_key": publicKey})
}
