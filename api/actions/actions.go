package actions

import (
	"encoding/json"
	"net/http"

	api_common "sisyphos/api/common"
	"sisyphos/lib/metadata"
	"sisyphos/lib/reqctx"
	"sisyphos/lib/utils"
	"sisyphos/models"
	gormrepo "sisyphos/repositories/gorm"
	"sisyphos/services"

	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"
)

var prep = func(r *http.Request) (*utils.Envelope, *services.ActionService) {
	env := reqctx.GetContext("envelope", r).(*utils.Envelope)
	actionRepo := gormrepo.NewActionRepo(reqctx.GetContext("db", r).(*gorm.DB), reqctx.GetContext("username", r).(string))
	actionService := services.NewActionService(actionRepo)
	return env, actionService
}

func Create(w http.ResponseWriter, r *http.Request) {
	env, actionService := prep(r)
	var m []models.Action
	err := json.NewDecoder(r.Body).Decode(&m)
	if err != nil {
		api_common.Render(w, r, env.SetError(err))
		return
	}
	actions, err := actionService.Create(m)
	if err != nil {
		api_common.Render(w, r, env.SetError(err))
		return
	}
	api_common.Render(w, r, env.SetResponse(actions))
}

func Update(w http.ResponseWriter, r *http.Request) {
	env, actionService := prep(r)
	name := chi.URLParam(r, "name")
	var m models.Action
	err := json.NewDecoder(r.Body).Decode(&m)
	if err != nil {
		api_common.Render(w, r, env.SetError(err))
		return
	}
	actions, err := actionService.Update(name, &m)
	if err != nil {
		api_common.Render(w, r, env.SetError(err))
		return
	}
	api_common.Render(w, r, env.SetResponse(actions))
}

func ReadOne(w http.ResponseWriter, r *http.Request) {
	env, actionService := prep(r)
	name := chi.URLParam(r, "name")
	actions, err := actionService.ReadByName(name)
	if err != nil {
		api_common.Render(w, r, env.SetError(err))
		return
	}
	api_common.Render(w, r, env.SetResponse(actions))
}

func ReadAll(w http.ResponseWriter, r *http.Request) {
	env, actionService := prep(r)
	md := reqctx.GetContext("metadata", r).(metadata.MetaData)
	var actions []models.Action
	var err error
	db := reqctx.GetContext("db", r).(*gorm.DB)
	if md.Filter != "" {
		// TODO clean up. no db here
		f := gormrepo.NewFilter(db, "actions")
		actions, err = actionService.ReadAllFiltered(md, f)
	} else {
		username := reqctx.GetContext("username", r).(string)
		permRepo := gormrepo.NewPermissionRepo(db, username)
		permService := services.NewPermissionService(permRepo)
		actionService.AddPermissionService(permService)

		actions, err = actionService.ReadAllExtendedPermission(username)
	}
	if err != nil {
		api_common.Render(w, r, env.SetError(err))
		return
	}
	api_common.Render(w, r, env.SetResponse(actions))
}

func Delete(w http.ResponseWriter, r *http.Request) {
}

func Execute(w http.ResponseWriter, r *http.Request) {
	env, actionService := prep(r)
	name := chi.URLParam(r, "name")
	extActions, err := actionService.ReadByName(name)
	if err != nil {
		api_common.Render(w, r, env.SetError(err))
		return
	}
	runRepo := gormrepo.NewRunRepo(
		reqctx.GetContext("db", r).(*gorm.DB),
		reqctx.GetContext("requestid", r).(string),
		reqctx.GetContext("username", r).(string),
	)
	runService := services.NewRunService(runRepo)
	actionService.AddRunService(runService)
	hostRepo := gormrepo.NewHostRepo(
		reqctx.GetContext("db", r).(*gorm.DB),
		reqctx.GetContext("username", r).(string),
	)
	hostService := services.NewHostService(hostRepo)
	actionService.AddHostService(hostService)
	runs, err := actionService.InitRun(extActions)
	if err != nil {
		api_common.Render(w, r, env.SetError(err))
		return
	}
	api_common.Render(w, r, env.SetResponse(runs))
}

func readRuns(w http.ResponseWriter, r *http.Request) {
	env, actionService := prep(r)
	name := chi.URLParam(r, "name")
	runs, err := actionService.ReadRuns(name)
	if err != nil {
		api_common.Render(w, r, env.SetError(err))
		return
	}
	api_common.Render(w, r, env.SetResponse(runs))
}
