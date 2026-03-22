package save

import (
	"log/slog"
	"net/http"
	resp "shortener/internal/lib/api/response"

	"shortener/internal/lib/logger/sl"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"shortener/internal/lib/random"
	"errors"
	"shortener/internal/storage"
)

type Request struct {
	URL   string `json:"url" validate:"required,url"`
	Alias string `json:"alias,omitempty"`
}

type Response struct {
	resp.Response
	Alias string `json:"alias,omitempty"`
}

// TODO: move to config if needed
const (
	aliasLength  = 6
	maxAliasRetries = 10
)

type URLSaver interface {
	SaveURL(urlToSave string, alias string) (int64, error)
	GetURL(alias string) (string, error)
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
		if err!= nil{
			log.Error("failed to decode request body", sl.Err(err))
			render.JSON(w,r,resp.Error("failed to decode request"))

			return
		}

		log.Info("request body decoded", slog.Any("request", req))

		if err:= validator.New().Struct(req); err!= nil{
			validateErr:= err.(validator.ValidationErrors)
			log.Error("invalid request", sl.Err(err))
			render.JSON(w,r,resp.ValidationError(validateErr))

			return
		}

		alias := req.Alias
		if alias == "" {
			for i := 0; i < maxAliasRetries; i++ {
				alias = random.NewRandomString(aliasLength)
				_, err := urlSaver.GetURL(alias)
				if errors.Is(err, storage.ErrURLNotFound) {
					break // алиас свободен
				}
				if err != nil {
					log.Error("failed to check alias", sl.Err(err))
					render.JSON(w, r, resp.Error("failed to check alias"))
					return
				}
				// алиас уже занят, перегенерируем
				alias = ""
			}
			if alias == "" {
				log.Error("failed to generate unique alias after retries")
				render.JSON(w, r, resp.Error("failed to generate unique alias"))
				return
			}
		}

		id, err:=urlSaver.SaveURL(req.URL, alias)
		if errors.Is(err, storage.ErrURLExists){
			log.Info("url already exists", slog.String("url", req.URL))
			render.JSON(w,r,resp.Error("url already exists"))
			return
		}
		if err!= nil{
			log.Error("failed to save url", sl.Err(err))
			render.JSON(w,r,resp.Error("failed to save url"))
			return
		}
		log.Info("url saved", slog.String("url", req.URL))
	}
}
