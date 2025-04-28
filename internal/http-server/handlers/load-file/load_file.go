package loadfile

import (
	"fmt"
	"io"
	"log/slog"
	mymiddleware "msu-logging-backend/internal/http-server/middleware"
	"msu-logging-backend/internal/lib/api/response"
	"msu-logging-backend/internal/services/audioservice"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-chi/render"
	"github.com/golang-jwt/jwt/v5"
)

func NewLoadFileHandler(log *slog.Logger, audioService *audioservice.AudioService) http.HandlerFunc {
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

		// Ограничиваем размер файла (например, 10 МБ)
		r.ParseMultipartForm(1 << 30)

		// Получаем файл из запроса
		file, handler, err := r.FormFile("audioFile")
		if err != nil {
			log.Error("Error Retrieving the File")
			return
		}
		defer file.Close()

		log.Info(fmt.Sprintf("Uploaded File: %+v\n", handler.Filename))
		log.Info(fmt.Sprintf("File Size: %+v\n", handler.Size))
		log.Info(fmt.Sprintf("MIME Header: %+v\n", handler.Header))

		// Создаём новый файл на сервере
		ext := filepath.Ext(handler.Filename)
		nameWithoutExt := strings.TrimSuffix(handler.Filename, ext)
		filename := fmt.Sprintf("%v_%v%v", nameWithoutExt, taskId, ext)

		dst, err := os.Create(filename)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer dst.Close()

		// Копируем содержимое загруженного файла в новый файл
		if _, err := io.Copy(dst, file); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		fmt.Fprintf(w, "Successfully Uploaded File\n")

		err = audioService.StartFileProcessing(taskId, filename)
		if err != nil {
			log.Error("Errpr in file processing")
		}

		render.JSON(w, r, response.Error("invalid token"))
	}
}
