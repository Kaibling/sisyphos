package apimiddleware

import (
	"context"
	"errors"
	"net/http"
	"strings"

	api_common "sisyphos/api/common"
	"sisyphos/lib/reqctx"
	"sisyphos/lib/utils"
	gormrepo "sisyphos/repositories/gorm"
	"sisyphos/services"

	"gorm.io/gorm"
)

func ValidateToken(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		env := reqctx.GetContext("envelope", r).(*utils.Envelope)
		if _, ok := r.Header["Authorization"]; !ok {
			api_common.Render(w, r, env.SetError(errors.New("authorization header missing")))
			return
		}
		if len(r.Header["Authorization"]) != 1 {
			api_common.Render(w, r, env.SetError(errors.New("authorization header missing")))
			return
		}
		authSlice := strings.Split(r.Header["Authorization"][0], " ")
		if len(authSlice) != 2 {
			api_common.Render(w, r, env.SetError(errors.New("authorization header invalid")))
			return
		}
		db := reqctx.GetContext("db", r).(*gorm.DB)
		userRepo := gormrepo.NewUserRepo(db, "unauthenticated")
		userService := services.NewUserService(userRepo)
		user, err := userService.ValidateToken(authSlice[1])
		if err != nil {
			api_common.Render(w, r, env.SetError(err))
			return
		}
		ctx := context.WithValue(r.Context(), reqctx.String("username"), user.Name)
		ctx = context.WithValue(ctx, reqctx.String("token"), authSlice[1])
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
