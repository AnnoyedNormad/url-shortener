package save

import (
	"errors"
	"fmt"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"log/slog"
	"net/http"
	"url-shortener/internal/config"
	"url-shortener/internal/lib/api/response"
	"url-shortener/internal/lib/logger/sl"
	"url-shortener/internal/lib/random"
	"url-shortener/internal/storage"
)

type Request struct {
	URL   string `json:"url" validate:"required,url"`
	Alias string `json:"alias,omitempty"`
}

type Response struct {
	response.Response
	Alias string `json:"alias,omitempty"`
}

type URLSaver interface {
	SaveURL(username, url, alias string) error
	AliasIsExist(alias string) (bool, error)
}

func New(log *slog.Logger, urlSaver URLSaver, cfg *config.Config) http.HandlerFunc {
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
			alias, err = generateRandomString(urlSaver, cfg)
			if err != nil {
				log.Error("can't generate random string", sl.Err(err))

				render.JSON(w, r, response.Error("internal error"))

				return
			}
		}

		user, _, _ := r.BasicAuth()
		err = urlSaver.SaveURL(user, req.URL, alias)
		if errors.Is(err, storage.ErrURLExists) {
			log.Info("url already exists", slog.String("url", req.URL), slog.String("err", err.Error()))

			render.JSON(w, r, response.Error("url already exists"))

			return
		}
		if err != nil {
			log.Error("failed to add url", slog.String("url", req.URL), slog.String("err", err.Error()))

			render.JSON(w, r, response.Error("failed to add url"))

			return
		}

		log.Info("url added", slog.String("user", user))

		render.JSON(w, r, Response{
			Response: response.OK(),
			Alias:    alias,
		})
	}
}

func generateRandomString(urlSaver URLSaver, cfg *config.Config) (string, error) {
	op := "generateRandomString"

	var alias = random.NewRandomString(cfg.AliasGenerator.Length, cfg.AliasGenerator.Alphabet)
	isExist, err := urlSaver.AliasIsExist(alias)
	var counter int
	for isExist {
		if counter == cfg.AliasGenerator.GenLimit {
			return "", fmt.Errorf("%s: %w", op, random.ErrAliasesAreOver)
		}
		if !errors.Is(err, storage.ErrURLExists) {
			return "", fmt.Errorf("%s: %w", op, err)
		}

		alias = random.NewRandomString(cfg.AliasGenerator.Length, cfg.AliasGenerator.Alphabet)
		isExist, err = urlSaver.AliasIsExist(alias)
		counter++
	}

	return alias, nil
}
