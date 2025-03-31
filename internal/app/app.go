package app

import (
	"log/slog"
	grpcapp "msu-logging-backend/internal/app/grpc"
	httpapp "msu-logging-backend/internal/app/http"
	minioapp "msu-logging-backend/internal/app/minio"
	rmqapp "msu-logging-backend/internal/app/rmq"
	wsapp "msu-logging-backend/internal/app/websocket"
	"msu-logging-backend/internal/config"
	"msu-logging-backend/internal/services/audioservice"
	"msu-logging-backend/internal/storage/mysql"
)

type App struct {
	GRPCSrv  *grpcapp.App
	WSSrv    *wsapp.App
	RMQSrv   *rmqapp.App
	MinioSrv *minioapp.App
	HTTPSrv  *httpapp.App
}

func New(
	log *slog.Logger,
	cfg *config.Config,
) *App {
	storage, err := mysql.New()
	if err != nil {
		panic(err)
	}

	app := &App{}

	app.MinioSrv = minioapp.New(log)
	app.RMQSrv = rmqapp.New(log, cfg)

	audio_service := audioservice.New(log, storage, app.RMQSrv, app.MinioSrv, cfg.MessageBroker.TranscribeQueue)
	app.GRPCSrv = grpcapp.New(log, cfg.GRPC.Port)
	app.WSSrv = wsapp.New(log, cfg.Websocket.Port, audio_service, cfg.Websocket.CertFile, cfg.Websocket.KeyFile)
	app.HTTPSrv = httpapp.New(log, cfg.HTTP.Address, storage)

	return app
}
