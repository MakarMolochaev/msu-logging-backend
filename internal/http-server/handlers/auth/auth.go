package auth

import (
	"context"
	"log/slog"
	"msu-logging-backend/internal/lib/api/response"
	jwtservice "msu-logging-backend/internal/services/jwt"
	"net/http"
	"time"

	"github.com/go-chi/render"
)

type Request struct {
}

type Response struct {
	response.Response
	Token     string    `json:"token"`
	ExpiresAt time.Time `json:"expires_at"`
}

type TaskStatusCreater interface {
	CreateNewTaskStatus(ctx context.Context) (int64, error)
}

func NewTokenHandler(log *slog.Logger, taskStatusCreater TaskStatusCreater, tokenTTL time.Duration) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.auth.NewTokenHandler"

		log = log.With(
			slog.String("op", op),
		)

		taskId, err := taskStatusCreater.CreateNewTaskStatus(context.Background())
		if err != nil {
			log.Error("Failed to save task_status in DB", slog.String("error", err.Error()))
			render.JSON(w, r, response.Error("Failed to save task_status in DB"))
			return
		}

		jwtService := jwtservice.New(log)
		tokenString, err := jwtService.GenerateToken(taskId, tokenTTL)

		if err != nil {
			log.Error("Failed to generate token", slog.String("error", err.Error()))
			render.JSON(w, r, response.Error("Failed to generate token"))
			return
		}

		render.JSON(w, r, Response{
			Response:  response.OK(),
			Token:     tokenString,
			ExpiresAt: time.Now().Add(tokenTTL),
		})
	}
}
