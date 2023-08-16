package groups

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"sisyphos/lib/reqctx"
	"sisyphos/lib/utils"
	"sisyphos/models"
	gormrepo "sisyphos/repositories/gorm"
	"sisyphos/services"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"gorm.io/gorm"
)

var prep = func(r *http.Request) (*utils.Envelope, *services.GroupService) {
	env := reqctx.GetContext("envelope", r).(*utils.Envelope)
	groupRepo := gormrepo.NewGroupRepo(reqctx.GetContext("db", r).(*gorm.DB))
	groupService := services.NewGroupService(groupRepo)
	return env, groupService
}

func Create(w http.ResponseWriter, r *http.Request) {
	env, groupService := prep(r)
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		render.Render(w, r, env.SetError(err))
		return
	}
	var m []models.Group
	err = json.Unmarshal(body, &m)
	if err != nil {
		render.Render(w, r, env.SetError(err))
		return
	}
	groups, err := groupService.Create(m)
	if err != nil {
		render.Render(w, r, env.SetError(err))
		return
	}
	render.Render(w, r, env.SetResponse(groups))
}

func Update(w http.ResponseWriter, r *http.Request) {
	env, groupService := prep(r)
	name := chi.URLParam(r, "name")
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		render.Render(w, r, env.SetError(err))
		return
	}
	var m *models.Group
	err = json.Unmarshal(body, &m)
	if err != nil {
		render.Render(w, r, env.SetError(err))
		return
	}
	groups, err := groupService.Update(name, m)
	if err != nil {
		render.Render(w, r, env.SetError(err))
		return
	}
	render.Render(w, r, env.SetResponse(groups))
}

func ReadOne(w http.ResponseWriter, r *http.Request) {
	env, groupService := prep(r)
	name := chi.URLParam(r, "name")
	groups, err := groupService.ReadByName(name)
	if err != nil {
		render.Render(w, r, env.SetError(err))
		return
	}
	render.Render(w, r, env.SetResponse(groups))
}

func ReadAll(w http.ResponseWriter, r *http.Request) {
	env, groupService := prep(r)
	groups, err := groupService.ReadAll()
	if err != nil {
		render.Render(w, r, env.SetError(err))
		return
	}
	render.Render(w, r, env.SetResponse(groups))
}

func Delete(w http.ResponseWriter, r *http.Request) {
}
