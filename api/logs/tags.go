package logs

import (
	"net/http"

	"sisyphos/lib/reqctx"
	"sisyphos/lib/utils"
	gormrepo "sisyphos/repositories/gorm"
	"sisyphos/services"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"gorm.io/gorm"
)

var prep = func(r *http.Request) (*utils.Envelope, *services.LogService) {
	env := reqctx.GetContext("envelope", r).(*utils.Envelope)
	logRepo := gormrepo.NewLogRepo(reqctx.GetContext("db", r).(*gorm.DB))
	logService := services.NewLogService(logRepo)
	return env, logService
}

func ReadOne(w http.ResponseWriter, r *http.Request) {
	env, logService := prep(r)
	id := chi.URLParam(r, "id")
	logs, err := logService.ReadByRequestID(id)
	if err != nil {
		render.Render(w, r, env.SetError(err))
		return
	}
	render.Render(w, r, env.SetResponse(logs))
}

func ReadAll(w http.ResponseWriter, r *http.Request) {
	env, logService := prep(r)
	logs, err := logService.ReadAll()
	if err != nil {
		render.Render(w, r, env.SetError(err))
		return
	}
	render.Render(w, r, env.SetResponse(logs))
}
