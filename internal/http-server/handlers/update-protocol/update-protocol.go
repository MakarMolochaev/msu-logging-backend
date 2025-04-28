package updateprotocol

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	minioapp "msu-logging-backend/internal/app/minio"
	mymiddleware "msu-logging-backend/internal/http-server/middleware"
	"msu-logging-backend/internal/lib/api/response"
	"net/http"
	"os"

	"github.com/go-chi/render"
	"github.com/golang-jwt/jwt/v5"
)

type RequestData struct {
	NewProtocol string `json:"new_protocol"`
}

type ProtocolUpdater interface {
	UpdateProtocolShortText(ctx context.Context, taskId int32, protocol string) (int64, error)
	UpdateProtocolFullText(ctx context.Context, taskId int32, full_text string) (int64, error)
}

func NewUpdateProtocolHandler(log *slog.Logger, protocolUpdater ProtocolUpdater, minioService *minioapp.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.auth.NewLoadFileHandler"

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

		var data RequestData
		err := json.NewDecoder(r.Body).Decode(&data)
		if err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}
		defer r.Body.Close()

		filename := fmt.Sprintf("protocol_%v.txt", taskId)
		file, err := os.Create(filename)
		if err != nil {
			log.Error("File creation error:", slog.String("error", err.Error()))
			http.Error(w, "Error in creating file", http.StatusBadRequest)
			return
		}
		_, err = file.Write([]byte(data.NewProtocol))
		if err != nil {
			log.Error("File writing error:", slog.String("error", err.Error()))
			http.Error(w, "Error in writing file", http.StatusBadRequest)
			return
		}

		link, err := minioService.UploadFile(filename, filename)
		if err != nil {
			log.Error("Minio upload error", slog.String("error", err.Error()))
			render.JSON(w, r, response.Error("Failed to save in MinIO"))
			return
		}

		log.Info("Audiofile uploaded to minio succesfully")

		protocolUpdater.UpdateProtocolShortText(context.Background(), taskId, link)

	}
}
