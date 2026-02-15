package handler

import (
	"errors"
	"io"
	"log/slog"
	"net/http"

	"github.com/ProgrammerBuffalo/url-shortener/internal/errs"
	"github.com/ProgrammerBuffalo/url-shortener/internal/service"
	"github.com/go-chi/render"
)

type Request struct {
	URL string `json:"url" validate:"required,url"`
}

type Response struct {
	Error      string `json:"error,omitempty"`
	StatusCode int    `json:"status_code,omitempty"`
	ShortURL   string `json:"short_url,omitempty"`
}

type CreateHandler struct {
	s      *service.UrlService
	logger *slog.Logger
}

func NewSaveHandler(s *service.UrlService, logger *slog.Logger) *CreateHandler {
	return &CreateHandler{s: s, logger: logger}
}

func (h *CreateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	const op = "handler.Create"
	h.logger.With("op", op)

	h.logger.Info("Begin creation url handler")
	var req Request

	if err := render.DecodeJSON(r.Body, &req); err != nil {
		render.Status(r, http.StatusBadRequest)
		if errors.Is(err, io.EOF) {
			render.JSON(w, r, Response{Error: "empty request", StatusCode: http.StatusBadRequest})
			return
		}
		render.JSON(w, r, Response{Error: "invalid request", StatusCode: http.StatusBadRequest})
		return
	}

	shortUrl, err := h.s.Create(r.Context(), req.URL)

	if err != nil {
		if errors.Is(err, errs.ErrDuplicateUrl) {
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, Response{Error: "duplicate long url", StatusCode: http.StatusBadRequest})
			return
		}
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, Response{Error: err.Error(), StatusCode: http.StatusBadRequest})
		return
	}

	h.logger.Info("Short url created", "shortUrl", "longUrl", shortUrl, req.URL)

	render.JSON(w, r, Response{ShortURL: shortUrl, StatusCode: http.StatusCreated})
}
