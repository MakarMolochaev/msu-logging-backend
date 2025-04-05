package valuation

import (
	"context"
	"log/slog"
	"msu-logging-backend/internal/lib/api/response"
	"net/http"

	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
)

type Request struct {
	Usability         int    `json:"usability"`
	ProcessingSpeed   int    `json:"processing_speed"`
	ProcessingQuality int    `json:"processing_quality"`
	ReuseService      bool   `json:"reuse_service"`
	Comment           string `json:"comment"`
}

type Response struct {
	response.Response
}

type ValuationSaver interface {
	SaveValuation(ctx context.Context, usability, processing_speed, processing_quality int, reuse_service bool, comment string) (int64, error)
}

func NewRateHandler(log *slog.Logger, valuationSaver ValuationSaver) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.valuation.NewRateHandler"

		log = log.With(
			slog.String("op", op),
		)

		var req Request

		err := render.DecodeJSON(r.Body, &req)
		if err != nil {
			log.Error("failed to decode request body", slog.String("error", err.Error()))

			render.JSON(w, r, response.Error("failed to decode request"))
			return
		}

		log.Info("request body decoded", slog.Any("request", req))

		if err := validator.New().Struct(req); err != nil {
			log.Error("invalid request", slog.String("error", err.Error()))

			render.JSON(w, r, response.Error("ivalid request"))

			return
		}

		_, err = valuationSaver.SaveValuation(context.Background(), req.Usability, req.ProcessingSpeed, req.ProcessingQuality, req.ReuseService, req.Comment)
		if err != nil {
			log.Error("error in saveing valudation", slog.String("error", err.Error()))
			return
		}

		log.Info("valuation saved")

		render.JSON(w, r, response.OK())
	}
}
