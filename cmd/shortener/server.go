package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/russianlagman/go-musthave-shortener/internal/app/handler/basic"
	"github.com/russianlagman/go-musthave-shortener/internal/app/handler/json"
	app "github.com/russianlagman/go-musthave-shortener/internal/app/middleware"
	"github.com/russianlagman/go-musthave-shortener/internal/app/service/store/sqlstore"
	"net/http"
)

func NewServer(config *Config, store *sqlstore.Store) *http.Server {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(app.SecureCookieAuth(config.SecretKey))
	r.Use(app.GzipResponseWriter)
	r.Use(app.GzipRequestReader)
	r.Get("/user/urls", json.UserDataHandler(store))
	r.With(app.ContentTypeJSON).Post("/api/shorten", json.WriteHandler(store))
	r.With(app.ContentTypeJSON).Post("/api/shorten/batch", json.BatchWriteHandler(store))
	r.With(app.ContentTypeJSON).Delete("/api/user/urls", json.BatchRemoveHandler(store))
	r.Get("/{id:[0-9a-z]+}", basic.ReadHandler(store))
	r.Post("/", basic.WriteHandler(store))
	r.Get("/ping", basic.PingHandler(store))

	return &http.Server{
		Addr:    config.ListenAddr,
		Handler: r,
	}
}
