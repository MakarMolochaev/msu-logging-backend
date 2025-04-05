package audiotask

import (
	"context"
	"log/slog"
	mymiddleware "msu-logging-backend/internal/http-server/middleware"
	"msu-logging-backend/internal/lib/api/response"
	"net/http"

	"github.com/go-chi/render"
	"github.com/golang-jwt/jwt"
)

type Request struct {
}

type Response struct {
	response.Response
	TaskStatus string `json:"task_status"`
}

type TaskStatusGetter interface {
	GetTaskStatusByID(ctx context.Context, id int64) (string, error)
}

func NewTaskStatusHandler(log *slog.Logger, taskStatusGetter TaskStatusGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.audiotask.NewTaskStatusHandler"

		log = log.With(
			slog.String("op", op),
		)

		claims, ok := r.Context().Value(mymiddleware.TokenClaimsKey).(jwt.MapClaims)
		if !ok {
			log.Error("failed to get JWT claims")
			render.JSON(w, r, response.Error("authentication failed"))
			return
		}

		taskId, ok := claims["taskId"].(int64)
		if !ok {
			log.Error("taskId claim not found or invalid")
			render.JSON(w, r, response.Error("invalid token"))
			return
		}

		taskStatus, err := taskStatusGetter.GetTaskStatusByID(r.Context(), taskId)
		if err != nil {
			log.Error("No task with this TaskId")
			render.JSON(w, r, response.Error("No task with this TaskId"))
			return
		}
		render.JSON(w, r, Response{
			Response:   response.OK(),
			TaskStatus: taskStatus,
		})

	}
}
