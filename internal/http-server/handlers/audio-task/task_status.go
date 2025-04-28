package audiotask

import (
	"context"
	"fmt"
	"log/slog"
	mymiddleware "msu-logging-backend/internal/http-server/middleware"
	"msu-logging-backend/internal/lib/api/response"
	"net/http"

	"github.com/go-chi/render"
	"github.com/golang-jwt/jwt/v5"
)

type Request struct {
}

type Response struct {
	response.Response
	TaskStatus    string `json:"task_status"`
	FullProtocol  string `json:"full_protocol"`
	ShortProtocol string `json:"short_protocol"`
}

type TaskStatusGetter interface {
	GetTaskStatusByID(ctx context.Context, id int32) (string, error)
}

type ProtocolGetter interface {
	GetProtocol(ctx context.Context, id int32) (string, string, error)
}

func NewTaskStatusHandler(log *slog.Logger, taskStatusGetter TaskStatusGetter, protocolGetter ProtocolGetter) http.HandlerFunc {
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

		taskClaim, ok := claims["taskId"]

		if !ok {
			log.Error("taskId claim not found or invalid")
			render.JSON(w, r, response.Error("invalid token"))
			return
		}

		taskId := int32(taskClaim.(float64))
		fmt.Println(taskId)
		taskStatus, err := taskStatusGetter.GetTaskStatusByID(r.Context(), taskId)
		if err != nil {
			log.Error("No task with this TaskId")
			render.JSON(w, r, response.Error("No task with this TaskId"))
			return
		}

		shortProtocolText, fullProtocolText, err := protocolGetter.GetProtocol(context.Background(), taskId)
		if err != nil {
			log.Error("No protocol with this TaskId")
			log.Info(string(taskId))
			render.JSON(w, r, response.Error("No protocol with this TaskId"))
			return
		}

		if taskStatus == "finished" {
			render.JSON(w, r, Response{
				Response:      response.OK(),
				TaskStatus:    taskStatus,
				FullProtocol:  fullProtocolText,
				ShortProtocol: shortProtocolText,
			})
		} else {
			render.JSON(w, r, Response{
				Response:   response.OK(),
				TaskStatus: taskStatus,
			})
		}
	}
}
