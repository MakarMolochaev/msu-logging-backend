package main

import (
	"log/slog"
	"msu-logging-backend/internal/app"
	"msu-logging-backend/internal/config"
	"os"
	"os/signal"
	"syscall"
)

func main() {

	cfg := config.MustLoad()

	log := setupLogger(cfg.Env)

	log.Info("Starting...")

	application := app.New(log, cfg)

	go application.GRPCSrv.MustRun()
	go application.WSSrv.MustRun()
	go application.RMQSrv.MustRun()
	go application.MinioSrv.MustRun()
	go application.HTTPSrv.MustRun()

	//shutdown

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	<-stop

	application.GRPCSrv.Stop()
	application.WSSrv.Stop()
	application.RMQSrv.Stop()
	log.Info("Application stopped")
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case "local":
		log = slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case "dev":
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case "prod":
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	}

	return log
}
