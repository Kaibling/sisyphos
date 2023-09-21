package apimiddleware

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"sisyphos/lib/reqctx"
	"sisyphos/models"
)

func ReadQueryParams(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		md := models.MetaData{}
		if val, ok := r.URL.Query()["filter"]; ok {
			if len(val) > 0 {
				md.Filter = val[0]
			}
		}
		if val, ok := r.URL.Query()["limit"]; ok {
			if len(val) > 0 {
				if l, err := strconv.Atoi(val[0]); err != nil {
					// TODO log error
					fmt.Println(err.Error())
				} else {
					md.Limit = l
				}
			}
		}
		if val, ok := r.URL.Query()["order"]; ok {
			if len(val) > 0 {
				md.Order = val[0]
			}
		}
		if val, ok := r.URL.Query()["sort_field"]; ok {
			if len(val) > 0 {
				md.SortField = val[0]
			}
		}
		// TODO set defaults for metadata
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
