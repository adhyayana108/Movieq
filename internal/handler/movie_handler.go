package handler

import (
	"errors"
	"net/http"

	"github.com/adhyayana108/movieq/internal/domain"
	"github.com/adhyayana108/movieq/internal/service"
)

type MovieHandler struct {
	svc service.MovieService
}

func NewMovieHandler(svc service.MovieService) *MovieHandler {
	return &MovieHandler{svc: svc}
}

//GET /movies

func (h *MovieHandler) GetMovies(w http.ResponseWriter, r *http.Request) {
	movies, err := h.svc.GetAll(r.Context())
	if err != nil {
		writeServiceError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, movies)
}

// GET /movies/{id}

func (h *MovieHandler) GetMovie(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	movie, err := h.svc.GetByID(r.Context(), id)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, movie)
}

// POST /movies

func (h *MovieHandler) CreateMovie(w http.ResponseWriter, r *http.Request) {
	var input domain.Movie
	if err := decodeJSON(r, &input); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	created, err := h.svc.Create(r.Context(), input)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, created)
}

// PUT /movies/{id}

func (h *MovieHandler) UpdateMovie(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	var input domain.Movie
	if err := decodeJSON(r, &input); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	updated, err := h.svc.Update(r.Context(), id, input)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, updated)
}

// DELETE /movies/{id}

func (h *MovieHandler) DeleteMovie(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	if err := h.svc.Delete(r.Context(), id); err != nil {
		writeServiceError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func writeServiceError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, domain.ErrNotFound):
		writeError(w, http.StatusNotFound, err.Error())
	case errors.Is(err, domain.ErrValidation):
		writeError(w, http.StatusBadRequest, err.Error())
	default:
		writeError(w, http.StatusInternalServerError, "internal server error")
	}
}
