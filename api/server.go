package api

import (
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	log "github.com/sirupsen/logrus"
	"gitlab.com/signmykey/signmykey/builtin/authenticator"
	"gitlab.com/signmykey/signmykey/builtin/principals"
	"gitlab.com/signmykey/signmykey/builtin/signer"
)

// Config represents the config of the API webserver.
type Config struct {
	VaultAddress string
	VaultToken   string
	VaultPath    string
	VaultRole    string

	TTL string

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

	log.Printf("signmykey server listen on 0.0.0.0:8080")
	return http.ListenAndServe(":8080", Router())
}

// Router returns *chi.Mux config
func Router() *chi.Mux {
	router := chi.NewRouter()
	router.Use(
		middleware.RequestID,
		middleware.RealIP,
		middleware.Logger,
		middleware.Recoverer,
		middleware.Timeout(30*time.Second),
	)

	router.Route("/v1", func(r chi.Router) {
		r.Get("/ping", pingHandler)
		r.Post("/sign", signHandler)
		r.Get("/ca", caHandler)
	})

	return router
}
