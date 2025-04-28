package audioservice

import (
	"context"
	"fmt"
	"log/slog"
	minioapp "msu-logging-backend/internal/app/minio"
	rmqapp "msu-logging-backend/internal/app/rmq"
	rabbitmodels "msu-logging-backend/internal/domain/models"
	"os"
)

type AudioService struct {
	log               *slog.Logger
	minio             *minioapp.App
	linkSaver         LinkSaver
	taskStatusSaver   TaskStatusSaver
	messageBroker     *rmqapp.App
	toTranscribeQueue string
	toProtocolQueue   string
}

type LinkSaver interface {
	SaveAudioFile(ctx context.Context, link string) (int64, error)
	UpdateProtocolShortText(ctx context.Context, taskId int32, protocol string) (int64, error)
	UpdateProtocolFullText(ctx context.Context, taskId int32, full_text string) (int64, error)
}

type TaskStatusSaver interface {
	UpdateTaskStatusByID(ctx context.Context, id int32, task_status string) error
}

func New(
	log *slog.Logger,
	linkSaver LinkSaver,
	taskStatusSaver TaskStatusSaver,
	messageBroker *rmqapp.App,
	minio *minioapp.App,
	toTranscribeQueue string,
	toProtocolQueue string,
) *AudioService {
	return &AudioService{
		log:               log,
		linkSaver:         linkSaver,
		taskStatusSaver:   taskStatusSaver,
		messageBroker:     messageBroker,
		minio:             minio,
		toTranscribeQueue: toTranscribeQueue,
		toProtocolQueue:   toProtocolQueue,
	}
}

func (a *AudioService) StartFileProcessing(taskId int32, filename string) error {
	const op = "audioservice.WhenWebsocketClosed"

	log := a.log.With(
		slog.String("op", op),
	)

	link, err := a.minio.UploadFile(filename, filename)
	if err != nil {
		a.log.Error("Minio upload error", slog.String("error", err.Error()))
		return fmt.Errorf("%s:Minio upload error: %w", op, err)
	}

	log.Info("Audiofile uploaded to minio succesfully")

	os.Remove(filename)

	_, err = a.linkSaver.SaveAudioFile(context.Background(), link)
	if err != nil {
		log.Error("MySQL save error", slog.String("error", err.Error()))
		return fmt.Errorf("%s: MySQL save error: %w", op, err)
	}

	log.Info("Audiofile uploaded to MySQL succesfully")

	transcribeRequestData := rabbitmodels.TranscribeRequest{
		TaskId:        taskId,
		AudioFileLink: link,
	}

	err = a.messageBroker.SendTranscribeRequest(a.toTranscribeQueue, transcribeRequestData)
	if err != nil {
		log.Error("RabbitMQ publish error", slog.String("error", err.Error()))
		return fmt.Errorf("%s: RabbitMQ publish error: %w", op, err)
	}

	err = a.taskStatusSaver.UpdateTaskStatusByID(context.Background(), taskId, "transcribing")

	if err != nil {
		a.log.Error("Error while updating the task", slog.String("error", err.Error()))
		return fmt.Errorf("%s: Error while updating the task: %w", op, err)
	}

	return nil
}

func (a *AudioService) WhenAudioTranscribed(taskId int32, transcribedText string) error {
	const op = "audioservice.WhenAudioTranscribed"

	log := a.log.With(
		slog.String("op", op),
	)

	transcribtionFilename := fmt.Sprintf("transcribed_%v.txt", taskId)

	file, err := os.Create(transcribtionFilename)
	if err != nil {
		log.Error("File creation error:", slog.String("error", err.Error()))
		return fmt.Errorf("%s:File creation error: %w", op, err)
	}
	_, err = file.Write([]byte(transcribedText))
	if err != nil {
		log.Error("File writing error:", slog.String("error", err.Error()))
		return fmt.Errorf("%s:File writing error: %w", op, err)
	}

	protocolLink, err := a.minio.UploadFile(transcribtionFilename, transcribtionFilename)
	if err != nil {
		a.log.Error("Minio upload error", slog.String("error", err.Error()))
		return fmt.Errorf("%s:Minio upload error: %w", op, err)
	}

	log.Info(fmt.Sprintf("Transcribtion #%v uploaded to minio succesfully", taskId))
	file.Close()
	os.Remove(transcribtionFilename)

	_, err = a.linkSaver.UpdateProtocolFullText(context.Background(), taskId, protocolLink)
	if err != nil {
		log.Error("MySQL save error", slog.String("error", err.Error()))
		return fmt.Errorf("%s: MySQL save error: %w", op, err)
	}

	protocolRequestData := rabbitmodels.ProtocolRequest{
		TaskId:          taskId,
		TranscribedText: transcribedText,
	}

	err = a.messageBroker.SendProtocolRequest(a.toProtocolQueue, protocolRequestData)
	if err != nil {
		log.Error("RabbitMQ publish error", slog.String("error", err.Error()))
		return fmt.Errorf("%s: RabbitMQ publish error: %w", op, err)
	}

	err = a.taskStatusSaver.UpdateTaskStatusByID(context.Background(), taskId, "making protocol")
	if err != nil {
		a.log.Error("Error while updating the task", slog.String("error", err.Error()))
		return fmt.Errorf("%s: Error while updating the task: %w", op, err)
	}

	return nil
}

func (a *AudioService) WhenProtocolIsReady(taskId int32, protocolText string) error {
	const op = "audioservice.WhenProtocolIsReady"

	log := a.log.With(
		slog.String("op", op),
	)

	protocolFilename := fmt.Sprintf("protocol_%v.txt", taskId)

	file, err := os.Create(protocolFilename)
	if err != nil {
		log.Error("File creation error:", slog.String("error", err.Error()))
		return fmt.Errorf("%s:File creation error: %w", op, err)
	}
	_, err = file.Write([]byte(protocolText))
	if err != nil {
		log.Error("File writing error:", slog.String("error", err.Error()))
		return fmt.Errorf("%s:File writing error: %w", op, err)
	}

	protocolLink, err := a.minio.UploadFile(protocolFilename, protocolFilename)
	if err != nil {
		a.log.Error("Minio upload error", slog.String("error", err.Error()))
		return fmt.Errorf("%s:Minio upload error: %w", op, err)
	}

	log.Info(fmt.Sprintf("Protocol #%v uploaded to minio succesfully", taskId))
	file.Close()
	os.Remove(protocolFilename)

	_, err = a.linkSaver.UpdateProtocolShortText(context.Background(), taskId, protocolLink)
	if err != nil {
		log.Error("MySQL save error", slog.String("error", err.Error()))
		return fmt.Errorf("%s: MySQL save error: %w", op, err)
	}

	log.Info("Protocol uploaded to MySQL succesfully")

	err = a.taskStatusSaver.UpdateTaskStatusByID(context.Background(), taskId, "finished")

	if err != nil {
		log.Error("Error while updating the task", slog.String("error", err.Error()))
		return fmt.Errorf("%s: Error while updating the task:%w", op, err)
	}

	return nil
}
