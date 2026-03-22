package redirect

import (
	"errors"
	"log/slog"
	"net/http"
	"shortener/internal/lib/logger/sl"
	"shortener/internal/storage"

	resp "shortener/internal/lib/api/response"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
)

type URLGetter interface {
	GetURL(alias string) (string, error)
}

func New(log *slog.Logger, urlGetter URLGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.redirect.New"
		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)
		alias := chi.URLParam(r, "alias")

		if alias == "" {
			log.Error("alias is empty")
			render.JSON(w, r, resp.Error("invalid request"))
			return
		}
		resURL, err := urlGetter.GetURL(alias)
		if errors.Is(err, storage.ErrURLNotFound) {
			log.Error("url not found", slog.String("alias", alias))
			render.JSON(w, r, resp.Error("url not found"))
			return
		}
		if err != nil {
			log.Error("failed to get url", sl.Err(err))
			render.JSON(w, r, resp.Error("internal error"))
			return
		}
		log.Info("url found", slog.String("url", resURL))
		http.Redirect(w, r, resURL, http.StatusTemporaryRedirect)
		return
	}
}
