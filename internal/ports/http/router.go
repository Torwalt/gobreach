package http

import (
	"gobreach/internal/domains/breach"
	"log"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

type httpRouter struct {
	Router  chi.Router
	service BreachServer
	logger  *log.Logger
}

type BreachServer interface {
	GetByEmail(email string) ([]breach.Breach, *breach.Error)
}

func NewRouter(s BreachServer, l *log.Logger) *httpRouter {
	r := chi.NewRouter()
	hr := &httpRouter{Router: r, service: s, logger: l}

	addMiddleware(hr)
	addRoutes(hr)
	return hr
}

func addMiddleware(r *httpRouter) {
	r.Router.Use(middleware.RequestID)
	r.Router.Use(middleware.RealIP)
	r.Router.Use(middleware.Logger)
	r.Router.Use(middleware.Recoverer)
}

func addRoutes(r *httpRouter) {
	r.Router.Get("/", indexHandler)
	bR := newBreachRouter(&r.service, r.logger)
	r.Router.Mount("/breach", bR)
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(http.StatusText(http.StatusOK)))
}
