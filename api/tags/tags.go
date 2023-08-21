package tags

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

var prep = func(r *http.Request) (*utils.Envelope, *services.TagService) {
	env := reqctx.GetContext("envelope", r).(*utils.Envelope)
	tagRepo := gormrepo.NewTagRepo(reqctx.GetContext("db", r).(*gorm.DB), reqctx.GetContext("username", r).(string))
	tagService := services.NewTagService(tagRepo)
	return env, tagService
}

func Create(w http.ResponseWriter, r *http.Request) {
	env, tagService := prep(r)
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		render.Render(w, r, env.SetError(err))
		return
	}
	var m []models.Tag
	err = json.Unmarshal(body, &m)
	if err != nil {
		render.Render(w, r, env.SetError(err))
		return
	}
	tags, err := tagService.Create(m)
	if err != nil {
		render.Render(w, r, env.SetError(err))
		return
	}
	render.Render(w, r, env.SetResponse(tags))
}

func Update(w http.ResponseWriter, r *http.Request) {
	env, tagService := prep(r)
	name := chi.URLParam(r, "name")
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		render.Render(w, r, env.SetError(err))
		return
	}
	var m map[string]interface{}
	err = json.Unmarshal(body, &m)
	if err != nil {
		render.Render(w, r, env.SetError(err))
		return
	}
	tags, err := tagService.Update(name, m)
	if err != nil {
		render.Render(w, r, env.SetError(err))
		return
	}
	render.Render(w, r, env.SetResponse(tags))
}

func ReadOne(w http.ResponseWriter, r *http.Request) {
	env, tagService := prep(r)
	name := chi.URLParam(r, "name")
	tags, err := tagService.ReadByName(name)
	if err != nil {
		render.Render(w, r, env.SetError(err))
		return
	}
	render.Render(w, r, env.SetResponse(tags))
}

func ReadAll(w http.ResponseWriter, r *http.Request) {
	env, tagService := prep(r)
	tags, err := tagService.ReadAll()
	if err != nil {
		render.Render(w, r, env.SetError(err))
		return
	}
	render.Render(w, r, env.SetResponse(tags))
}

func Delete(w http.ResponseWriter, r *http.Request) {
}
