package api

import (
	"net/http"

	"github.com/go-chi/render"
)

func pingHandler(w http.ResponseWriter, r *http.Request) {
	render.JSON(w, r, map[string]string{"message": "pong"})
}
