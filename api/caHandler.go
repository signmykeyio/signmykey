package api

import (
	"net/http"

	"github.com/go-chi/render"
)

func caHandler(w http.ResponseWriter, r *http.Request) {
	publicKey, err := config.Signer.ReadCA()
	if err != nil {
		render.Status(r, 500)
		render.JSON(w, r, map[string]string{"errors": "error getting CA certificate"})
		return
	}

	render.JSON(w, r, map[string]string{"public_key": publicKey})
}
