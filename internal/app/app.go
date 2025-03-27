package app

import (
	"log/slog"
	grpcapp "msu-logging-backend/internal/app/grpc"
	"time"
)

type App struct {
	GRPCSrv *grpcapp.App
}

func New(
	log *slog.Logger,
	grpcPort int,
	tokenTTL time.Duration,
) *App {

	grpcApp := grpcapp.New(log, grpcPort)
	return &App{GRPCSrv: grpcApp}
}
