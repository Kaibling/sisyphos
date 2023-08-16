package logs

import (
	"sisyphos/lib/apimiddleware"

	"github.com/go-chi/chi/v5"
)

func AddRoutes() chi.Router {
	r := chi.NewRouter()
	r.Route("/", func(r chi.Router) {
		r.Use(apimiddleware.ValidateToken)
		r.Get("/", ReadAll)
		r.Get("/{id}", ReadOne)
	})

	return r
}
