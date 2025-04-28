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

type Response struct {
	response.Response
	Token     string    `json:"token"`
	ExpiresAt time.Time `json:"expires_at"`
}

type TaskStatusCreater interface {
	CreateNewTaskStatus(ctx context.Context) (int32, error)
	CreateNewProtocol(ctx context.Context, task_id int32) error
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

		err = taskStatusCreater.CreateNewProtocol(context.Background(), taskId)
		if err != nil {
			log.Error("Failed to save protocol placeholder in DB", slog.String("error", err.Error()))
			render.JSON(w, r, response.Error("Failed to save protocol placeholder in DB"))
			return
		}

		//w.Header().Set("Authorization", "Bearer "+tokenString)
		http.SetCookie(w, &http.Cookie{
			Name:     "jwt_token",
			Value:    tokenString,
			Path:     "/",
			Domain:   "",
			Secure:   false,
			HttpOnly: false,
			SameSite: http.SameSiteLaxMode,
			Expires:  time.Now().Add(tokenTTL),
		})

		render.JSON(w, r, Response{
			Response:  response.OK(),
			Token:     tokenString,
			ExpiresAt: time.Now().Add(tokenTTL),
		})
	}
}
