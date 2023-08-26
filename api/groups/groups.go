package groups

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

var prep = func(r *http.Request) (*utils.Envelope, *services.GroupService) {
	env := reqctx.GetContext("envelope", r).(*utils.Envelope)
	groupRepo := gormrepo.NewGroupRepo(reqctx.GetContext("db", r).(*gorm.DB), reqctx.GetContext("username", r).(string))
	groupService := services.NewGroupService(groupRepo)
	return env, groupService
}

func Create(w http.ResponseWriter, r *http.Request) {
	env, groupService := prep(r)
	body, err := io.ReadAll(r.Body)
	if err != nil {
		api_common.Render(w, r, env.SetError(err))
		return
	}
	var m []models.Group
	err = json.Unmarshal(body, &m)
	if err != nil {
		api_common.Render(w, r, env.SetError(err))
		return
	}
	groups, err := groupService.Create(m)
	if err != nil {
		api_common.Render(w, r, env.SetError(err))
		return
	}
	api_common.Render(w, r, env.SetResponse(groups))
}

func Update(w http.ResponseWriter, r *http.Request) {
	env, groupService := prep(r)
	name := chi.URLParam(r, "name")
	body, err := io.ReadAll(r.Body)
	if err != nil {
		api_common.Render(w, r, env.SetError(err))
		return
	}
	var m *models.Group
	err = json.Unmarshal(body, &m)
	if err != nil {
		api_common.Render(w, r, env.SetError(err))
		return
	}
	groups, err := groupService.Update(name, m)
	if err != nil {
		api_common.Render(w, r, env.SetError(err))
		return
	}
	api_common.Render(w, r, env.SetResponse(groups))
}

func ReadOne(w http.ResponseWriter, r *http.Request) {
	env, groupService := prep(r)
	name := chi.URLParam(r, "name")
	groups, err := groupService.ReadByName(name)
	if err != nil {
		api_common.Render(w, r, env.SetError(err))
		return
	}
	api_common.Render(w, r, env.SetResponse(groups))
}

func ReadAll(w http.ResponseWriter, r *http.Request) {
	env, groupService := prep(r)
	groups, err := groupService.ReadAll()
	if err != nil {
		api_common.Render(w, r, env.SetError(err))
		return
	}
	api_common.Render(w, r, env.SetResponse(groups))
}

func Delete(w http.ResponseWriter, r *http.Request) {
}
