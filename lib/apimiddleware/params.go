package apimiddleware

import (
	"bytes"
	"context"
	"io"
	"net/http"

	"sisyphos/lib/metadata"
	"sisyphos/lib/reqctx"
)

func ReadQueryParams(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		md := metadata.MetaData{}
		if val, ok := r.URL.Query()["filter"]; ok {
			if len(val) > 0 {
				md.Filter = val[0]
			}
		}
		ctx := context.WithValue(r.Context(), reqctx.String("metadata"), md)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func ReadBody(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		b, err := io.ReadAll(r.Body)
		if err != nil {
			return
		}
		r.Body = io.NopCloser(bytes.NewBuffer(b))
		ctx := context.WithValue(r.Context(), reqctx.String("bytebody"), b)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
