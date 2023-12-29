package save

import (
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"log/slog"
	"net/http"
	"url-shortener/internal/lib/api/response"
	"url-shortener/internal/lib/logger/sl"
	"url-shortener/internal/storage"
)

type Request struct {
	URL   string `json:"url" validate:"required,url"`
	Alias string `json:"alias,omitempty"`
}

type Response struct {
	response.Response
	URLS []storage.AliasUrl `json:"urls,omitempty"`
}

type AliasURLGetter interface {
	GetUserAliasURL(user string) ([]storage.AliasUrl, error)
}

func New(log *slog.Logger, aliasURLGetter AliasURLGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.url.geturls.New"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		user, _, _ := r.BasicAuth()
		urls, err := aliasURLGetter.GetUserAliasURL(user)
		if err != nil {
			log.Info("failed to get url", sl.Err(err))

			render.JSON(w, r, response.Error("internal error"))

			return
		}
		log.Info("got urls")

		render.JSON(w, r, Response{
			Response: response.OK(),
			URLS:     urls,
		})
	}
}
