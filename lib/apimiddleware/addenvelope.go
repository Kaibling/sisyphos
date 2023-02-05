package apimiddleware

import (
	"context"
	"net/http"

	"sisyphos/lib/utils"
)

func AddEnvelope(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var reqID string
		if val, ok := r.Header["X-Request-Id"]; ok {
			reqID = val[0]
		} else {
			reqID = utils.NewULID().String()
		}

		env := utils.NewEnvelope()
		env.RequestID = reqID
		ctx := context.WithValue(r.Context(), "envelope", env)
		ctx = context.WithValue(ctx, "requestid", reqID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
