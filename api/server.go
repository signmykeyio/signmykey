package api

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/signmykeyio/signmykey/builtin/authenticator"
	"github.com/signmykeyio/signmykey/builtin/principals"
	"github.com/signmykeyio/signmykey/builtin/signer"
	log "github.com/sirupsen/logrus"
)

// Config represents the config of the API webserver.
type Config struct {
	Addr       string
	TLSDisable bool
	TLSCert    string
	TLSKey     string

	Auth   authenticator.Authenticator
	Princs principals.Principals
	Signer signer.Signer
}

var (
	config Config
)

// Serve the API webserver and register all handlers
func Serve(startconfig Config) error {
	config = startconfig

	// Config logging
	formatter := &log.TextFormatter{
	//FullTimestamp: true,
	}
	log.SetFormatter(formatter)

	if config.TLSDisable {
		log.Printf("signmykey server listen on http://%s", config.Addr)
		return http.ListenAndServe(config.Addr, Router())
	}

	if _, err := os.Stat(config.TLSCert); os.IsNotExist(err) {
		return fmt.Errorf("Cert file %s doesn't exist", err)
	}

	if _, err := os.Stat(config.TLSKey); os.IsNotExist(err) {
		return fmt.Errorf("Key file %s doesn't exist", err)
	}

	log.Printf("signmykey server listen on https://%s", config.Addr)
	return http.ListenAndServeTLS(config.Addr, config.TLSCert, config.TLSKey, Router())
}

// Router returns *chi.Mux config
func Router() *chi.Mux {
	router := chi.NewRouter()
	router.Use(
		middleware.RequestID,
		middleware.RealIP,
		middleware.Logger,
		middleware.Recoverer,
		middleware.CloseNotify,
		middleware.Timeout(15*time.Second),
	)

	router.Route("/v1", func(r chi.Router) {
		r.Get("/ping", pingHandler)
		r.Post("/sign", signHandler)
		r.Get("/ca", caHandler)
	})

	return router
}
