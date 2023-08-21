package authentication

import (
	"encoding/json"
	"net/http"

	"sisyphos/lib/reqctx"
	"sisyphos/lib/utils"
	"sisyphos/models"
	gormrepo "sisyphos/repositories/gorm"
	"sisyphos/services"

	"github.com/go-chi/render"
	"gorm.io/gorm"
)

var prep = func(r *http.Request) (*utils.Envelope, *services.UserService) {
	env := reqctx.GetContext("envelope", r).(*utils.Envelope)
	userRepo := gormrepo.NewUserRepo(reqctx.GetContext("db", r).(*gorm.DB), "unauthenticated")
	userService := services.NewUserService(userRepo)
	return env, userService
}

func login(w http.ResponseWriter, r *http.Request) {
	env, userService := prep(r)
	var m models.Authentication
	if err := json.NewDecoder(r.Body).Decode(&m); err != nil {
		render.Render(w, r, env.SetError(err))
		return
	}
	user, err := userService.Authenticate(m)
	if err != nil {
		render.Render(w, r, env.SetError(err))
		return
	}
	render.Render(w, r, env.SetResponse(user))
}
