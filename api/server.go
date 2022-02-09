package api

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/signmykeyio/signmykey/builtin/authenticator"
	"github.com/signmykeyio/signmykey/builtin/principals"
	"github.com/signmykeyio/signmykey/builtin/signer"
	"github.com/sirupsen/logrus"
)

// Config represents the config of the API webserver.
type Config struct {
	Addr       string
	TLSDisable bool
	TLSCert    string
	TLSKey     string

	Logger *logrus.Logger

	Auth   authenticator.Authenticator
	Princs []principals.Principals
	Signer signer.Signer
}

type contextKey string

var (
	config Config
)

// Serve the API webserver and register all handlers
func Serve(startconfig Config) {
	config = startconfig

	// Config logging
	logger := config.Logger

	if config.TLSDisable {
		logger.WithField("ctx", "api").Warn("Running signmykey server with TLS disabled is strongly discouraged!")
		logger.WithField("ctx", "api").Infof("Signmykey server listen on http://%s", config.Addr)
		err := http.ListenAndServe(config.Addr, Router(logger))
		if err != nil {
			logger.WithField("ctx", "api").WithError(err).Error("Serving HTTP")
		}
		return
	}

	if _, err := os.Stat(config.TLSCert); os.IsNotExist(err) {
		logger.WithField("ctx", "api").WithError(err).Error("Load TLS certificate")
		return
	}

	if _, err := os.Stat(config.TLSKey); os.IsNotExist(err) {
		logger.WithField("ctx", "api").WithError(err).Error("Load TLS key")
		return
	}

	tlsCfg := &tls.Config{
		MinVersion:               tls.VersionTLS12,
		CurvePreferences:         []tls.CurveID{tls.CurveP521, tls.CurveP384, tls.CurveP256},
		PreferServerCipherSuites: true,
		CipherSuites: []uint16{
			// TLS 1.3
			tls.TLS_AES_128_GCM_SHA256,
			tls.TLS_AES_256_GCM_SHA384,
			tls.TLS_CHACHA20_POLY1305_SHA256,
			// TLS 1.2, support will be remove from future version
			tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
		},
	}

	// Create standard logger from logrus for http.Server internal ErrorLog
	lw := logger.WithField("ctx", "api").WriterLevel(logrus.ErrorLevel)
	defer lw.Close()

	srv := &http.Server{
		Addr:         config.Addr,
		Handler:      Router(logger),
		TLSConfig:    tlsCfg,
		TLSNextProto: make(map[string]func(*http.Server, *tls.Conn, http.Handler)),
		ErrorLog:     log.New(lw, "", 0),
	}

	logger.WithField("ctx", "api").Infof("Signmykey server listen on https://%s", config.Addr)
	err := srv.ListenAndServeTLS(config.TLSCert, config.TLSKey)
	if err != nil {
		logger.WithField("ctx", "api").WithError(err).Error("Serving HTTP request")
		return
	}
}

// Router returns *chi.Mux config
func Router(logger *logrus.Logger) *chi.Mux {
	router := chi.NewRouter()
	router.Use(
		middleware.RequestID,
		middleware.RealIP,
		Logger(logger),
		middleware.Recoverer,
		middleware.Timeout(15*time.Second),
	)

	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "<title>Signmykey</title>")
		fmt.Fprintf(w, "Welcome to <a href=https://signmykey.io/>Signmykey</a> service !")
	})

	router.Route("/v1", func(r chi.Router) {
		r.Get("/ping", pingHandler)
		r.Post("/sign", signHandler)
		r.Get("/ca", caHandler)
	})

	return router
}
