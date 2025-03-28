package audioservice

import (
	"context"
	"fmt"
	"log/slog"
	rmqapp "msu-logging-backend/internal/app/rmq"
)

type AudioService struct {
	log                *slog.Logger
	audioFileLinkSaver AudioFileLinkSaver
	messageBroker      *rmqapp.App
}

type AudioFileLinkSaver interface {
	SaveAudioFile(ctx context.Context, link string) (int64, error)
}

func New(
	log *slog.Logger,
	audioFileLinkSaver AudioFileLinkSaver,
	messageBroker *rmqapp.App,
) *AudioService {
	return &AudioService{
		log:                log,
		audioFileLinkSaver: audioFileLinkSaver,
		messageBroker:      messageBroker,
	}
}

func (a *AudioService) WebsocketClosed(msg []byte) {
	fmt.Println("Websocket closed, audio chunk:", len(msg), "bytes")
}
