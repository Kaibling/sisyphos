package apimiddleware

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"sisyphos/lib/utils"
	gormrepo "sisyphos/repositories/gorm"
	"sisyphos/services"

	"github.com/go-chi/render"
	"gorm.io/gorm"
)

func ValidateToken(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		env := utils.GetContext("envelope", r).(*utils.Envelope)
		if _, ok := r.Header["Authorization"]; !ok {
			render.Render(w, r, env.SetError(errors.New("authorization header missing")))
			return
		}
		if len(r.Header["Authorization"]) != 1 {
			render.Render(w, r, env.SetError(errors.New("authorization header missing")))
			return
		}
		authSlice := strings.Split(r.Header["Authorization"][0], " ")
		if len(authSlice) != 2 {
			render.Render(w, r, env.SetError(errors.New("authorization header invalid")))
			return
		}
		db := utils.GetContext("db", r).(*gorm.DB)
		userRepo := gormrepo.NewUserRepo(db)
		userService := services.NewUserService(userRepo)
		user, err := userService.ValidateToken(authSlice[1])
		if err != nil {
			render.Render(w, r, env.SetError(err))
			return
		}
		ctx := context.WithValue(r.Context(), "username", user.Name)
		ctx = context.WithValue(ctx, "token", authSlice[1])
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
