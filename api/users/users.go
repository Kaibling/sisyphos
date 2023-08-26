package users

import (
	"encoding/json"
	"io"
	"net/http"

	api_common "sisyphos/api/common"
	"sisyphos/lib/reqctx"
	"sisyphos/lib/utils"
	"sisyphos/models"
	gormrepo "sisyphos/repositories/gorm"
	"sisyphos/services"

	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"
)

var prep = func(r *http.Request) (*utils.Envelope, *services.UserService) {
	env := reqctx.GetContext("envelope", r).(*utils.Envelope)
	userRepo := gormrepo.NewUserRepo(reqctx.GetContext("db", r).(*gorm.DB), reqctx.GetContext("username", r).(string))
	userService := services.NewUserService(userRepo)
	return env, userService
}

func Create(w http.ResponseWriter, r *http.Request) {
	env, userService := prep(r)
	body, err := io.ReadAll(r.Body)
	if err != nil {
		api_common.Render(w, r, env.SetError(err))
		return
	}
	var m []models.User
	if err = json.Unmarshal(body, &m); err != nil {
		api_common.Render(w, r, env.SetError(err))
		return
	}
	if err := models.UserArrayValidate(m); err != nil {
		api_common.Render(w, r, env.SetError(err))
		return
	}

	users, err := userService.Create(m)
	if err != nil {
		api_common.Render(w, r, env.SetError(err))
		return
	}
	api_common.Render(w, r, env.SetResponse(users))
}

func Update(w http.ResponseWriter, r *http.Request) {
	env, userService := prep(r)
	name := chi.URLParam(r, "name")
	body, err := io.ReadAll(r.Body)
	if err != nil {
		api_common.Render(w, r, env.SetError(err))
		return
	}
	var m map[string]interface{}
	err = json.Unmarshal(body, &m)
	if err != nil {
		api_common.Render(w, r, env.SetError(err))
		return
	}
	users, err := userService.Update(name, m)
	if err != nil {
		api_common.Render(w, r, env.SetError(err))
		return
	}
	api_common.Render(w, r, env.SetResponse(users))
}

func ReadOne(w http.ResponseWriter, r *http.Request) {
	env, userService := prep(r)
	name := chi.URLParam(r, "name")
	users, err := userService.ReadByName(name)
	if err != nil {
		api_common.Render(w, r, env.SetError(err))
		return
	}
	api_common.Render(w, r, env.SetResponse(users))
}

func ReadAll(w http.ResponseWriter, r *http.Request) {
	env, userService := prep(r)
	permRepo := gormrepo.NewPermissionRepo(reqctx.GetContext("db", r).(*gorm.DB), reqctx.GetContext("username", r).(string))
	permService := services.NewPermissionService(permRepo)
	userService.AddPermissionService(permService)
	users, err := userService.ReadAllPermission(reqctx.GetContext("username", r).(string))
	// users, err := userService.ReadAll()
	if err != nil {
		api_common.Render(w, r, env.SetError(err))
		return
	}
	api_common.Render(w, r, env.SetResponse(users))
}

func Delete(w http.ResponseWriter, r *http.Request) {
}
