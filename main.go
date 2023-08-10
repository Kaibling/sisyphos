package main

import (
	"context"
	"fmt"
	"net/http"

	"sisyphos/api"
	"sisyphos/lib/apimiddleware"
	"sisyphos/lib/config"
	"sisyphos/lib/log"
	"sisyphos/models"

	gormrepo "sisyphos/repositories/gorm"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
)

func main() {
	config.Init()

	l, err := log.New()
	if err != nil {
		fmt.Println("logger creation failed: ", err.Error())
		return
	}
	ctx := context.WithValue(context.Background(), models.String("logger"), l)
	dbconfig := gormrepo.DBConfig{
		User:     config.Config.DBUser,
		Port:     config.Config.DBPort,
		Password: config.Config.DBPassword,
		Host:     config.Config.DBHost,
		Database: config.Config.DBDatabase,
	}
	db, err := gormrepo.InitDatabase(dbconfig, l)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	r := chi.NewRouter()

	r.Use(injectContextData("db", db))
	r.Use(injectContextData("logger", l))
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{"*"},
		// Access-Control-Allow-Origin
		AllowedMethods: []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{"*"},
		// AllowedHeaders:     []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		AllowCredentials: false,
		MaxAge:           300,
	}))
	r.Use(middleware.Logger)
	r.Use(apimiddleware.AddEnvelope)
	r.Use(apimiddleware.ReadQueryParams)

	r.Mount("/", api.Routes())

	//displayRoutes(r)

	//fmt.Println("listening on :3000")
	log.Info(ctx, "listening on :3000")
	//log.Info(ctx, fmt.Sprintf("listening on %s", listeningStr))
	err = http.ListenAndServe(":3000", r)
	if err != nil {
		fmt.Println(err.Error())
	}
}

func injectContextData(key string, data interface{}) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			ctx := context.WithValue(r.Context(), key, data)
			next.ServeHTTP(w, r.WithContext(ctx))
		}
		return http.HandlerFunc(fn)
	}
}

func displayRoutes(r *chi.Mux) {
	walkFunc := func(method string, route string, handler http.Handler, middlewares ...func(http.Handler) http.Handler) error {
		fmt.Printf("%s %s\n", method, route)
		return nil
	}
	if err := chi.Walk(r, walkFunc); err != nil {
		fmt.Printf("Logging err: %s\n", err.Error())
	}
}
