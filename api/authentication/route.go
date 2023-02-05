package authentication

import (
	"github.com/go-chi/chi/v5"
)

func AddRoutes() chi.Router {
	r := chi.NewRouter()
	r.Post("/login", login)
	return r
}
