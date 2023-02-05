package apimiddleware

import (
	"context"
	"net/http"

	"sisyphos/lib/metadata"
)

func ReadQueryParams(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		md := metadata.MetaData{}
		if val, ok := r.URL.Query()["filter"]; ok {
			if len(val) > 0 {
				md.Filter = val[0]
			}
		}
		ctx := context.WithValue(r.Context(), "metadata", md)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
