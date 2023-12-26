package save

import (
	"errors"
	"fmt"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"log/slog"
	"net/http"
	"url-shortener/internal/lib/api/response"
	"url-shortener/internal/lib/logger/sl"
	"url-shortener/internal/lib/random"
	"url-shortener/internal/storage"
)

// TODO: load alias length from config
const aliasLength = 7

type Request struct {
	URL   string `json:"url" validate:"required,url"`
	Alias string `json:"alias,omitempty"`
}

type Response struct {
	response.Response
	Alias string `json:"alias,omitempty"`
}

type URLSaver interface {
	SaveURL(url string, alias string) error
	URLIsExist(alias string) (bool, error)
}

func New(log *slog.Logger, urlSaver URLSaver) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.url.save.New"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var req Request

		err := render.DecodeJSON(r.Body, &req)
		if err != nil {
			log.Error("failed to decode request")

			render.JSON(w, r, response.Error("failed to decode request"))

			return
		}

		log.Info("request body decoded", slog.Any("request", req))

		if err = validator.New().Struct(req); err != nil {
			log.Error("invalid request", sl.Err(err))

			render.JSON(w, r, response.Error("invalid request"))

			return
		}

		alias := req.Alias
		if alias == "" {
			alias, err = generateRandomString(urlSaver)
			if err != nil {
				log.Error("can't generate random string", sl.Err(err))

				render.JSON(w, r, response.Error("can't generate random string"))

				return
			}
		}

		err = urlSaver.SaveURL(req.URL, alias)
		if errors.Is(err, storage.ErrURLExists) {
			log.Info("url already exists", slog.String("url", req.URL))

			render.JSON(w, r, response.Error("url already exists"))

			return
		}
		if err != nil {
			log.Error("failed to add url", slog.String("url", req.URL))

			render.JSON(w, r, response.Error("failed to add url"))

			return
		}

		log.Info("url added")

		render.JSON(w, r, Response{
			Response: response.OK(),
			Alias:    alias,
		})
	}
}

func generateRandomString(urlSaver URLSaver) (string, error) {
	op := "generateRandomString"
	var alias = random.NewRandomString(aliasLength)

	isExist, err := urlSaver.URLIsExist(alias)
	for !isExist {
		if err != nil {
			return "", fmt.Errorf("%s: %w", op, err)
		}

		isExist, err = urlSaver.URLIsExist(alias)
	}

	return alias, nil
}
