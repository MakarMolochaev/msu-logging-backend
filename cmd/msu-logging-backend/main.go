package main

import (
	"fmt"
	"log"
	"log/slog"
	"msu-logging-backend/internal/app"
	"msu-logging-backend/internal/config"
	"os"
	"os/signal"
	"syscall"

	"github.com/gorilla/websocket"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {

	cfg := config.MustLoad()

	log := setupLogger(cfg.Env)

	log.Info("Starting...")

	application := app.New(log, cfg, handleWebSocketConnection)

	go application.GRPCSrv.MustRun()
	go application.WSSrv.MustRun()

	//shutdown

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	<-stop

	application.GRPCSrv.Stop()
	application.WSSrv.Stop()
	log.Info("Application stopped")
}

func handleWebSocketConnection(conn *websocket.Conn) {
	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			log.Println("Client disconnected:", err)
			return
		}
		// Обработка аудио-данных
		fmt.Println("Received audio chunk:", len(msg), "bytes")
	}
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal:
		log = slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envDev:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envProd:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	}

	return log
}
