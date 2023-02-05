package users

import (
	"sisyphos/lib/apimiddleware"

	"github.com/go-chi/chi/v5"
)

func AddRoutes() chi.Router {
	r := chi.NewRouter()
	r.Route("/", func(r chi.Router) {
		r.Use(apimiddleware.ValidateToken)
		r.Post("/", Create)
		r.Get("/", ReadAll)
		r.Get("/{name}", ReadOne)
		r.Delete("/{name}", Delete)
		r.Patch("/{name}", Update)
	})

	return r
}
