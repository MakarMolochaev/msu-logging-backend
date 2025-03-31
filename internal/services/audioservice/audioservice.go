package audioservice

import (
	"context"
	"fmt"
	"log/slog"
	minioapp "msu-logging-backend/internal/app/minio"
	rmqapp "msu-logging-backend/internal/app/rmq"
)

type AudioService struct {
	log                *slog.Logger
	minio              *minioapp.App
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
	minio *minioapp.App,
) *AudioService {
	return &AudioService{
		log:                log,
		audioFileLinkSaver: audioFileLinkSaver,
		messageBroker:      messageBroker,
		minio:              minio,
	}
}

func (a *AudioService) WhenWebsocketClosed(filename string) error {
	const op = "audioservice.WhenWebsocketClosed"

	/*
		log := a.log.With(
			slog.String("op", op),
		)
	*/

	link, err := a.minio.UploadFile(filename, filename)
	if err != nil {
		a.log.Error("Read error", slog.String("error", err.Error()))
		return fmt.Errorf("file upload error: %w", err)
	}

	//os.Remove(filename)
	//log.Info("Audiofile uploaded to minio succesfully")

	_, err = a.audioFileLinkSaver.SaveAudioFile(context.Background(), link)
	if err != nil {
		return err
	}

	//log.Info("Audiofile uploaded to MySQL succesfully")
	return nil
}
