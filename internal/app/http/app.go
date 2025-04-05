package httpapp

import (
	"log/slog"
	"msu-logging-backend/internal/config"
	audiotask "msu-logging-backend/internal/http-server/handlers/audio-task"
	"msu-logging-backend/internal/http-server/handlers/auth"
	"msu-logging-backend/internal/http-server/handlers/valuation"
	mymiddleware "msu-logging-backend/internal/http-server/middleware"
	"msu-logging-backend/internal/storage/mysql"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
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
	config *config.Config,
) *App {

	router := chi.NewRouter()
	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://example.com", "https://*.example.com"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))
	router.Use(middleware.RequestID)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)
	router.Post("/valuation", valuation.NewRateHandler(log, storage))
	router.Get("/token", auth.NewTokenHandler(log, storage, config.HTTP.TokenTTL))

	router.Group(func(r chi.Router) {
		r.Use(mymiddleware.JWTVerifier(os.Getenv("JWT_SECRET")))
		r.Get("/taskstatus", audiotask.NewTaskStatusHandler(log, storage))
	})

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
