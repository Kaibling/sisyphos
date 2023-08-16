package apimiddleware

import (
	"context"
	"net/http"

	"sisyphos/lib/reqctx"
	"sisyphos/lib/utils"
	gormrepo "sisyphos/repositories/gorm"
	"sisyphos/services"

	"gorm.io/gorm"
)

func AddEnvelope(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var reqID string
		if val, ok := r.Header["X-Request-Id"]; ok {
			reqID = val[0]
		} else {
			reqID = utils.NewULID().String()
		}

		db := reqctx.GetContext("db", r).(*gorm.DB)
		lr := gormrepo.NewLogRepo(db)
		ls := services.NewLogService(lr)
		env := utils.NewEnvelope(ls)
		env.RequestID = reqID
		ctx := context.WithValue(r.Context(), reqctx.String("envelope"), env)
		ctx = context.WithValue(ctx, reqctx.String("requestid"), reqID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
