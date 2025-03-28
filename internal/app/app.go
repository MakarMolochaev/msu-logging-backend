package app

import (
	"log/slog"
	grpcapp "msu-logging-backend/internal/app/grpc"
	rmqapp "msu-logging-backend/internal/app/rmq"
	wsapp "msu-logging-backend/internal/app/websocket"
	"msu-logging-backend/internal/config"
	"msu-logging-backend/internal/services/audioservice"
	"msu-logging-backend/internal/storage/mysql"
)

type App struct {
	GRPCSrv *grpcapp.App
	WSSrv   *wsapp.App
	RMQSrv  *rmqapp.App
}

func New(
	log *slog.Logger,
	cfg *config.Config,
) *App {

	/*
		storage, err := mysql.New()
		if err != nil {
			panic(err)
		}
	*/
	storage, err := mysql.New()
	if err != nil {
		panic(err)
	}

	var rmqApp *rmqapp.App
	var grpcApp *grpcapp.App
	var wsApp *wsapp.App

	audio_service := audioservice.New(log, storage, rmqApp)

	rmqApp = rmqapp.New(log, cfg)
	grpcApp = grpcapp.New(log, cfg.GRPC.Port)
	wsApp = wsapp.New(log, cfg.Websocket.Port, audio_service, cfg.Websocket.CertFile, cfg.Websocket.KeyFile)
	return &App{GRPCSrv: grpcApp, WSSrv: wsApp, RMQSrv: rmqApp}
}
