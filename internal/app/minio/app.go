package minioapp

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type App struct {
	log         *slog.Logger
	client      *minio.Client
	bucket_name string
}

func New(
	log *slog.Logger,
) *App {
	return &App{
		log:         log,
		bucket_name: os.Getenv("MINIO_BUCKET_NAME"),
	}
}

func (a *App) Run() error {
	const op = "minioapp.Run"

	log := a.log.With(slog.String("op", op))

	client, err := minio.New(os.Getenv("MINIO_ENDPOINT"), &minio.Options{
		Creds:  credentials.NewStaticV4(os.Getenv("MINIO_USER"), os.Getenv("MINIO_PASSWORD"), ""),
		Secure: false,
	})
	if err != nil {
		log.Error("Failed to connect to MinIO")
		return fmt.Errorf("%s: ошибка при создании MinIO клиента: %w", op, err)
	}

	a.client = client
	log.Info("Minio is ready")
	return nil
}

func (a *App) MustRun() {
	if err := a.Run(); err != nil {
		panic(err)
	}
}

func (a *App) UploadFile(objectName, filePath string) (string, error) {
	const op = "minioapp.UploadFile"

	log := a.log.With(
		slog.String("op", op),
	)
	log.Info("Uploading file to minio...")

	_, err := a.client.FPutObject(context.Background(), a.bucket_name, objectName, filePath, minio.PutObjectOptions{})
	if err != nil {
		return "", fmt.Errorf("%s: Ошибка при загрузке файла: %w", op, err)
	}

	log.Info("Файл %s успешно загружен в бакет %s\n", objectName, a.bucket_name)

	link, err := a.client.PresignedGetObject(context.Background(), a.bucket_name, objectName, time.Hour, nil)
	if err != nil {
		return "", fmt.Errorf("%s: Ошибка при получении временной ссылки на файл: %w", op, err)
	}
	return link.String(), nil
}
