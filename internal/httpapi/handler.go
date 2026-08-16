package httpapi

import (
	"embed"
	"encoding/json"
	"errors"
	"io"
	"io/fs"
	"net/http"
	"strings"

	"course-trial/internal/trial"
	"course-trial/web"
	"github.com/go-chi/chi/v5"
)

type Handler struct {
	service *trial.Service
}

func New(service *trial.Service) http.Handler {
	h := &Handler{service: service}
	r := chi.NewRouter()
	r.Get("/api/pages/{slug}", h.page)
	r.Post("/api/pages/{slug}/visits", h.start)
	r.Post("/api/visits/{visitID}/next", h.next)
	r.Get("/api/admin/pages/{slug}", h.adminPage)
	r.Put("/api/admin/pages/{slug}", h.updatePage)
	assets, _ := fs.Sub(web.Assets, "src")
	r.Handle("/*", http.FileServer(http.FS(assets)))
	return r
}

func (h *Handler) adminPage(w http.ResponseWriter, r *http.Request) {
	page, err := h.service.AdminPage(chi.URLParam(r, "slug"))
	if err != nil {
		h.writeError(w, err)
		return
	}
	h.writeJSON(w, http.StatusOK, page)
}

func (h *Handler) updatePage(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	var update trial.PageUpdate
	if err := decoder.Decode(&update); err != nil {
		h.writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid_request", "message": "invalid page request"})
		return
	}
	if err := decoder.Decode(&struct{}{}); !errors.Is(err, io.EOF) {
		h.writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid_request", "message": "invalid page request"})
		return
	}
	page, err := h.service.UpdatePage(chi.URLParam(r, "slug"), update)
	if err != nil {
		h.writeError(w, err)
		return
	}
	h.writeJSON(w, http.StatusOK, page)
}

func (h *Handler) page(w http.ResponseWriter, r *http.Request) {
	view, err := h.service.Page(chi.URLParam(r, "slug"))
	if err != nil {
		h.writeError(w, err)
		return
	}
	h.writeJSON(w, http.StatusOK, view)
}

func (h *Handler) start(w http.ResponseWriter, r *http.Request) {
	started, err := h.service.Start(chi.URLParam(r, "slug"))
	if err != nil {
		if errors.Is(err, trial.ErrAccessLimit) {
			view, viewErr := h.service.Page(chi.URLParam(r, "slug"))
			if viewErr == nil {
				h.writeJSON(w, http.StatusConflict, map[string]string{"error": "access_limit", "message": strings.TrimSpace(err.Error()), "closing_copy": view.ClosingCopy})
				return
			}
		}
		h.writeError(w, err)
		return
	}
	h.writeJSON(w, http.StatusCreated, started)
}

func (h *Handler) next(w http.ResponseWriter, r *http.Request) {
	step, err := h.service.Next(chi.URLParam(r, "visitID"))
	if err != nil {
		h.writeError(w, err)
		return
	}
	h.writeJSON(w, http.StatusOK, step)
}

func (h *Handler) writeError(w http.ResponseWriter, err error) {
	status := http.StatusInternalServerError
	code := "internal_error"
	switch {
	case errors.Is(err, trial.ErrPageNotFound):
		status, code = http.StatusNotFound, "page_not_found"
	case errors.Is(err, trial.ErrVisitNotFound):
		status, code = http.StatusNotFound, "visit_not_found"
	case errors.Is(err, trial.ErrPageInactive):
		status, code = http.StatusGone, "page_inactive"
	case errors.Is(err, trial.ErrAccessLimit):
		status, code = http.StatusConflict, "access_limit"
	case errors.Is(err, trial.ErrInvalidPage):
		status, code = http.StatusUnprocessableEntity, "invalid_page"
	}
	h.writeJSON(w, status, map[string]string{"error": code, "message": strings.TrimSpace(err.Error())})
}

func (h *Handler) writeJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(value)
}

var _ embed.FS
