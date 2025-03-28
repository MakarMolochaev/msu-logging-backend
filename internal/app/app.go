package app

import (
	"log/slog"
	grpcapp "msu-logging-backend/internal/app/grpc"
	wsapp "msu-logging-backend/internal/app/websocket"
	"msu-logging-backend/internal/config"

	"github.com/gorilla/websocket"
)

type App struct {
	GRPCSrv *grpcapp.App
	WSSrv   *wsapp.App
}

func New(
	log *slog.Logger,
	cfg *config.Config,
	wsHandler func(*websocket.Conn),
) *App {

	/*
		storage, err := mysql.New()
		if err != nil {
			panic(err)
		}
	*/
	grpcApp := grpcapp.New(log, cfg.GRPC.Port)
	wsApp := wsapp.New(log, cfg.Websocket.Port, wsHandler, cfg.Websocket.CertFile, cfg.Websocket.KeyFile)
	return &App{GRPCSrv: grpcApp, WSSrv: wsApp}
}
