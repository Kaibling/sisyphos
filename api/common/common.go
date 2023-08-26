package common

import (
	"net/http"

	"sisyphos/lib/log"

	"github.com/go-chi/render"
)

func Render(w http.ResponseWriter, r *http.Request, v render.Renderer) {
	if err := render.Render(w, r, v); err != nil {
		log.Error(r.Context(), err)
	}
}
