package actions

import (
	"encoding/json"
	"net/http"

	"sisyphos/lib/metadata"
	"sisyphos/lib/utils"
	"sisyphos/models"
	gormrepo "sisyphos/repositories/gorm"
	"sisyphos/services"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"gorm.io/gorm"
)

var prep = func(r *http.Request) (*utils.Envelope, *services.ActionService) {
	env := utils.GetContext("envelope", r).(*utils.Envelope)
	actionRepo := gormrepo.NewActionRepo(utils.GetContext("db", r).(*gorm.DB))
	actionService := services.NewActionService(actionRepo)
	return env, actionService
}

func Create(w http.ResponseWriter, r *http.Request) {
	env, actionService := prep(r)
	var m []models.Action
	err := json.NewDecoder(r.Body).Decode(&m)
	// body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		render.Render(w, r, env.SetError(err))
		return
	}

	// err = json.Unmarshal(body, &m)
	// if err != nil {
	// 	render.Render(w, r, env.SetError(err))
	// 	return
	// }
	actions, err := actionService.Create(m)
	if err != nil {
		render.Render(w, r, env.SetError(err))
		return
	}
	render.Render(w, r, env.SetResponse(actions))
}

func Update(w http.ResponseWriter, r *http.Request) {
	env, actionService := prep(r)
	name := chi.URLParam(r, "name")
	var m models.Action //map[string]interface{}
	err := json.NewDecoder(r.Body).Decode(&m)
	if err != nil {
		render.Render(w, r, env.SetError(err))
		return
	}
	actions, err := actionService.Update(name, &m)
	if err != nil {
		render.Render(w, r, env.SetError(err))
		return
	}
	render.Render(w, r, env.SetResponse(actions))
}

func ReadOne(w http.ResponseWriter, r *http.Request) {
	env, actionService := prep(r)
	name := chi.URLParam(r, "name")
	actions, err := actionService.ReadByName(name)
	if err != nil {
		render.Render(w, r, env.SetError(err))
		return
	}
	render.Render(w, r, env.SetResponse(actions))
}

func ReadAll(w http.ResponseWriter, r *http.Request) {
	env, actionService := prep(r)
	md := utils.GetContext("metadata", r).(metadata.MetaData)
	var actions []models.Action
	var err error
	db := utils.GetContext("db", r).(*gorm.DB)
	if md.Filter != "" {
		// TODO clean up. no db here

		f := gormrepo.NewFilter(db, "actions")
		actions, err = actionService.ReadAllFiltered(md, f)
	} else {
		// actions, err = actionService.ReadAllExtended()
		permRepo := gormrepo.NewPermissionRepo(db)
		permService := services.NewPermissionService(permRepo)
		actionService.AddPermissionService(permService)
		username := utils.GetContext("username", r).(string)
		actions, err = actionService.ReadAllExtendedPermission(username)
	}
	if err != nil {
		render.Render(w, r, env.SetError(err))
		return
	}
	render.Render(w, r, env.SetResponse(actions))
}

func Delete(w http.ResponseWriter, r *http.Request) {
}

func Execute(w http.ResponseWriter, r *http.Request) {
	env, actionService := prep(r)
	name := chi.URLParam(r, "name")
	extActions, err := actionService.ReadExt(name)
	if err != nil {
		render.Render(w, r, env.SetError(err))
		return
	}
	runRepo := gormrepo.NewRunRepo(
		utils.GetContext("db", r).(*gorm.DB),
		utils.GetContext("requestid", r).(string),
		utils.GetContext("username", r).(string),
	)
	runService := services.NewRunService(runRepo)
	actionService.AddRunService(runService)
	hostRepo := gormrepo.NewHostRepo(
		utils.GetContext("db", r).(*gorm.DB),
	)
	hostService := services.NewHostService(hostRepo)
	actionService.AddHostService(hostService)
	runs, err := actionService.InitRun(extActions)
	if err != nil {
		render.Render(w, r, env.SetError(err))
		return
	}
	render.Render(w, r, env.SetResponse(runs))
}

func readRuns(w http.ResponseWriter, r *http.Request) {
	env, actionService := prep(r)
	name := chi.URLParam(r, "name")
	runs, err := actionService.ReadRuns(name)
	if err != nil {
		render.Render(w, r, env.SetError(err))
		return
	}
	render.Render(w, r, env.SetResponse(runs))
}
