package handler

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/adhyayana108/movieq/internal/domain"
)

type fakeService struct {
	movies    []domain.Movie
	movie     domain.Movie
	err       error
	createArg domain.Movie
}

func (f *fakeService) GetAll(_ context.Context) ([]domain.Movie, error) { return f.movies, f.err }
func (f *fakeService) GetByID(_ context.Context, _ string) (domain.Movie, error) {
	return f.movie, f.err
}
func (f *fakeService) Create(_ context.Context, m domain.Movie) (domain.Movie, error) {
	f.createArg = m
	if f.err != nil {
		return domain.Movie{}, f.err
	}
	m.ID = "new-id"
	return m, nil
}
func (f *fakeService) Update(_ context.Context, _ string, m domain.Movie) (domain.Movie, error) {
	return m, f.err
}
func (f *fakeService) Delete(_ context.Context, _ string) error { return f.err }

func TestGetMovies_OK(t *testing.T) {
	svc := &fakeService{movies: []domain.Movie{{ID: "1", Title: "Arrival", ISBN: "123"}}}
	h := NewMovieHandler(svc)

	req := httptest.NewRequest(http.MethodGet, "/movies", nil)
	rec := httptest.NewRecorder()

	h.GetMovies(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("got status %d, want %d", rec.Code, http.StatusOK)
	}
	var got []domain.Movie
	if err := json.Unmarshal(rec.Body.Bytes(), &got); err != nil {
		t.Fatalf("response body did not decode as []domain.Movie: %v", err)
	}
	if len(got) != 1 || got[0].Title != "Arrival" {
		t.Errorf("got %+v, want one movie titled Arrival", got)
	}
}

func TestGetMovie_NotFound_Returns404(t *testing.T) {
	svc := &fakeService{err: domain.ErrNotFound}
	h := NewMovieHandler(svc)

	req := httptest.NewRequest(http.MethodGet, "/movies/999", nil)
	req.SetPathValue("id", "999")
	rec := httptest.NewRecorder()

	h.GetMovie(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Errorf("got status %d, want %d", rec.Code, http.StatusNotFound)
	}
}

func TestCreateMovie_ValidBody_Returns201(t *testing.T) {
	svc := &fakeService{}
	h := NewMovieHandler(svc)

	body := `{"isbn":"123","title":"Dune","director":{"firstname":"Denis","lastname":"Villeneuve"}}`
	req := httptest.NewRequest(http.MethodPost, "/movies", strings.NewReader(body))
	rec := httptest.NewRecorder()

	h.CreateMovie(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("got status %d, want %d; body: %s", rec.Code, http.StatusCreated, rec.Body.String())
	}
	if svc.createArg.Title != "Dune" {
		t.Errorf("service received title %q, want %q", svc.createArg.Title, "Dune")
	}
}

func TestCreateMovie_MalformedJSON_Returns400(t *testing.T) {
	svc := &fakeService{}
	h := NewMovieHandler(svc)

	req := httptest.NewRequest(http.MethodPost, "/movies", strings.NewReader(`{"title": "Dune"`)) // truncated
	rec := httptest.NewRecorder()

	h.CreateMovie(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("got status %d, want %d", rec.Code, http.StatusBadRequest)
	}
}

func TestCreateMovie_UnknownField_Returns400(t *testing.T) {
	svc := &fakeService{}
	h := NewMovieHandler(svc)

	req := httptest.NewRequest(http.MethodPost, "/movies", strings.NewReader(`{"titel":"Dune","isbn":"1"}`))
	rec := httptest.NewRecorder()

	h.CreateMovie(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("got status %d, want %d", rec.Code, http.StatusBadRequest)
	}
}

func TestCreateMovie_UnknownServiceError_Returns500(t *testing.T) {
	svc := &fakeService{err: errors.New("validation failed: title is required")}
	h := NewMovieHandler(svc)

	req := httptest.NewRequest(http.MethodPost, "/movies", strings.NewReader(`{"isbn":"1","title":"x"}`))
	rec := httptest.NewRecorder()

	h.CreateMovie(rec, req)

	if rec.Code != http.StatusInternalServerError {
		t.Errorf("got status %d, want %d", rec.Code, http.StatusInternalServerError)
	}
}

func TestDeleteMovie_Returns204(t *testing.T) {
	svc := &fakeService{}
	h := NewMovieHandler(svc)

	req := httptest.NewRequest(http.MethodDelete, "/movies/1", nil)
	req.SetPathValue("id", "1")
	rec := httptest.NewRecorder()

	h.DeleteMovie(rec, req)

	if rec.Code != http.StatusNoContent {
		t.Errorf("got status %d, want %d", rec.Code, http.StatusNoContent)
	}
	if rec.Body.Len() != 0 {
		t.Errorf("got body %q, want empty body for 204", rec.Body.String())
	}
}