package httpapp

import (
	"log/slog"
	"msu-logging-backend/internal/http-server/handlers/valuation"
	"msu-logging-backend/internal/storage/mysql"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type App struct {
	log        *slog.Logger
	HTTPServer *http.Server
	address    string
}

func New(
	log *slog.Logger,
	address string,
	storage *mysql.Storage,
) *App {

	router := chi.NewRouter()
	router.Use(middleware.RequestID)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)
	router.Post("/valuation", valuation.New(log, storage))

	HTTPServer := &http.Server{
		Addr:    address,
		Handler: router,
	}

	return &App{
		log:        log,
		HTTPServer: HTTPServer,
		address:    address,
	}
}

func (a *App) Run() error {
	const op = "httpapp.Run"

	log := a.log.With(
		slog.String("op", op),
	)

	a.log.Info("starting HTTP server", slog.String("adress", a.address))

	if err := a.HTTPServer.ListenAndServe(); err != nil {
		log.Error("failed to start server")
	}

	return nil
}

func (a *App) MustRun() {
	if err := a.Run(); err != nil {
		panic(err)
	}
}

func (a *App) Stop() {
	const op = "httpapp.stop"

	a.log.With(slog.String("op", op)).
		Info("Stopping HTTP server")
}
