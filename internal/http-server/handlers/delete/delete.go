package delete

import (
	"errors"
	"log/slog"
	"net/http"
	resp "shortener/internal/lib/api/response"
	"shortener/internal/lib/logger/sl"
	"shortener/internal/storage"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

type URLDeleter interface {
	DeleteURL(alias string) error
}

func New(log *slog.Logger, urlDeleter URLDeleter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		alias := chi.URLParam(r, "alias")
		if alias == "" {
			render.JSON(w, r, resp.Error("invalid request"))
			return
		}
		err := urlDeleter.DeleteURL(alias)
		if errors.Is(err, storage.ErrURLNotFound) {
			render.JSON(w, r, resp.Error("url not found"))
			return
		}
		if err != nil {
			log.Error("failed to delete url", sl.Err(err))
			render.JSON(w, r, resp.Error("internal error"))
			return
		}
		render.JSON(w, r, resp.OK())
	}
}
