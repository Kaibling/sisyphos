package api

import (
	"sisyphos/api/actions"
	"sisyphos/api/authentication"
	"sisyphos/api/groups"
	"sisyphos/api/hosts"
	"sisyphos/api/logs"
	"sisyphos/api/tags"
	"sisyphos/api/users"

	"github.com/go-chi/chi/v5"
)

func Routes() *chi.Mux {
	r := chi.NewRouter()
	r.Group(func(r chi.Router) {
		r.Mount("/actions", actions.AddRoutes())
		r.Mount("/hosts", hosts.AddRoutes())
		r.Mount("/logs", logs.AddRoutes())
		r.Mount("/tags", tags.AddRoutes())
		r.Mount("/users", users.AddRoutes())
		r.Mount("/groups", groups.AddRoutes())
		r.Mount("/authentication", authentication.AddRoutes())
	})
	return r
}
