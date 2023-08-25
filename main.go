package main

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"sisyphos/api"
	"sisyphos/lib/apimiddleware"
	"sisyphos/lib/cluster"
	"sisyphos/lib/cluster/repos/postgres"
	"sisyphos/lib/config"
	"sisyphos/lib/log"
	"sisyphos/lib/reqctx"

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
	l.AddDefaultField("component", "api")
	ctx := context.WithValue(context.Background(), reqctx.String("logger"), l)
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
	r.Use(apimiddleware.ReadBody)

	r.Mount("/", api.Routes())

	if config.Config.ClusterEnabled {
		cl := l.Copy()
		cl.AddDefaultField("component", "cluster")
		cfg := cluster.ClusterConfig{
			StartHook:    func() { fmt.Println("master fuck jeah") },
			StopHook:     func() {},
			HeatBeatRate: time.Duration(config.Config.ClusterHeatBeatRate) * time.Millisecond,
			Log:          cl,
		}

		be := postgres.New(
			postgres.PostgresConfig{
				User:     config.Config.DBUser,
				Port:     config.Config.DBPort,
				Password: config.Config.DBPassword,
				Host:     config.Config.DBHost,
				Database: config.Config.DBDatabase,
			}, cl)

		if err := be.Connect(); err != nil {
			log.Error(ctx, err)
			return
		}

		c, err := cluster.New(cfg, be)
		if err != nil {
			log.Error(ctx, err)
			return
		}

		if err := be.AddEmptyLock(c.ID()); err != nil {
			log.Error(ctx, err)
			return
		}
		go func() {
			if err := c.Run(); err != nil {
				log.Error(ctx, err)
				return
			}
		}()
	}

	listeningStr := fmt.Sprintf("%s:%s", config.Config.BindingIP, config.Config.BindingPort)
	log.Info(ctx, fmt.Sprintf("listening on %s", listeningStr))
	err = http.ListenAndServe(listeningStr, r)
	if err != nil {
		fmt.Println(err.Error())
	}
}

func injectContextData(key reqctx.String, data interface{}) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			ctx := context.WithValue(r.Context(), key, data)
			next.ServeHTTP(w, r.WithContext(ctx))
		}
		return http.HandlerFunc(fn)
	}
}
