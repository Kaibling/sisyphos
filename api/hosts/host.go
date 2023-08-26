package hosts

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

var prep = func(r *http.Request) (*utils.Envelope, *services.HostService) {
	env := reqctx.GetContext("envelope", r).(*utils.Envelope)
	hostRepo := gormrepo.NewHostRepo(reqctx.GetContext("db", r).(*gorm.DB), reqctx.GetContext("username", r).(string))
	hostService := services.NewHostService(hostRepo)
	return env, hostService
}

func Create(w http.ResponseWriter, r *http.Request) {
	env, hostService := prep(r)
	body, err := io.ReadAll(r.Body)
	if err != nil {
		api_common.Render(w, r, env.SetError(err))
		return
	}
	var m []*models.Host
	err = json.Unmarshal(body, &m)
	if err != nil {
		api_common.Render(w, r, env.SetError(err))
		return
	}
	hosts, err := hostService.Create(m)
	if err != nil {
		api_common.Render(w, r, env.SetError(err))
		return
	}
	api_common.Render(w, r, env.SetResponse(hosts))
}

func Update(w http.ResponseWriter, r *http.Request) {
	env, hostService := prep(r)
	name := chi.URLParam(r, "name")
	body, err := io.ReadAll(r.Body)
	if err != nil {
		api_common.Render(w, r, env.SetError(err))
		return
	}
	var m *models.Host
	err = json.Unmarshal(body, &m)
	if err != nil {
		api_common.Render(w, r, env.SetError(err))
		return
	}
	hosts, err := hostService.Update(name, m)
	if err != nil {
		api_common.Render(w, r, env.SetError(err))
		return
	}
	api_common.Render(w, r, env.SetResponse(hosts))
}

func ReadOne(w http.ResponseWriter, r *http.Request) {
	env, hostService := prep(r)
	name := chi.URLParam(r, "name")
	hosts, err := hostService.ReadByName(name)
	if err != nil {
		api_common.Render(w, r, env.SetError(err))
		return
	}
	api_common.Render(w, r, env.SetResponse(hosts))
}

func ReadAll(w http.ResponseWriter, r *http.Request) {
	env, hostService := prep(r)
	hosts, err := hostService.ReadAll()
	if err != nil {
		api_common.Render(w, r, env.SetError(err))
		return
	}
	api_common.Render(w, r, env.SetResponse(hosts))
}

func Delete(w http.ResponseWriter, r *http.Request) {
}

func testConnection(w http.ResponseWriter, r *http.Request) {
	env, hostService := prep(r)
	name := chi.URLParam(r, "name")
	err := hostService.TestConnection(name)
	if err != nil {
		api_common.Render(w, r, env.SetError(err))
		return
	}
	api_common.Render(w, r, env.SetResponse("ok"))
}
