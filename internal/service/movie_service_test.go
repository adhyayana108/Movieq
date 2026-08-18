package service

import (
	"context"
	"errors"
	"testing"

	"github.com/adhyayana108/movieq/internal/domain"
)

type fakeRepository struct {
	movies    map[string]domain.Movie
	forceErr  error
	callCount int
}

func newFakeRepository() *fakeRepository {
	return &fakeRepository{movies: make(map[string]domain.Movie)}
}

func (f *fakeRepository) GetAll(_ context.Context) ([]domain.Movie, error) {
	f.callCount++
	if f.forceErr != nil {
		return nil, f.forceErr
	}
	out := make([]domain.Movie, 0, len(f.movies))
	for _, m := range f.movies {
		out = append(out, m)
	}
	return out, nil
}

func (f *fakeRepository) GetByID(_ context.Context, id string) (domain.Movie, error) {
	f.callCount++
	if f.forceErr != nil {
		return domain.Movie{}, f.forceErr
	}
	m, ok := f.movies[id]
	if !ok {
		return domain.Movie{}, domain.ErrNotFound
	}
	return m, nil
}

func (f *fakeRepository) Create(_ context.Context, m domain.Movie) (domain.Movie, error) {
	f.callCount++
	m.ID = "fake-id"
	f.movies[m.ID] = m
	return m, nil
}

func (f *fakeRepository) Update(_ context.Context, id string, m domain.Movie) (domain.Movie, error) {
	f.callCount++
	if _, ok := f.movies[id]; !ok {
		return domain.Movie{}, domain.ErrNotFound
	}
	m.ID = id
	f.movies[id] = m
	return m, nil
}

func (f *fakeRepository) Delete(_ context.Context, id string) error {
	f.callCount++
	if _, ok := f.movies[id]; !ok {
		return domain.ErrNotFound
	}
	delete(f.movies, id)
	return nil
}

func TestMovieService_Create_RejectsMissingTitle(t *testing.T) {
	repo := newFakeRepository()
	svc := NewMovieService(repo)

	_, err := svc.Create(context.Background(), domain.Movie{ISBN: "123"})
	if !errors.Is(err, domain.ErrValidation) {
		t.Fatalf("got err %v, want domain.ErrValidation", err)
	}
	if repo.callCount != 0 {
		t.Errorf("repository was called %d times; validation should fail before storage is touched", repo.callCount)
	}
}

func TestMovieService_Create_RejectsMissingISBN(t *testing.T) {
	repo := newFakeRepository()
	svc := NewMovieService(repo)

	_, err := svc.Create(context.Background(), domain.Movie{Title: "Dune"})
	if !errors.Is(err, domain.ErrValidation) {
		t.Fatalf("got err %v, want domain.ErrValidation", err)
	}
}

func TestMovieService_Create_AcceptsValidMovie(t *testing.T) {
	repo := newFakeRepository()
	svc := NewMovieService(repo)

	created, err := svc.Create(context.Background(), domain.Movie{
		Title: "Dune",
		ISBN:  "999",
		Director: &domain.Director{
			FirstName: "Denis",
			LastName:  "Villeneuve",
		},
	})
	if err != nil {
		t.Fatalf("Create returned unexpected error: %v", err)
	}
	if created.ID == "" {
		t.Error("Create did not assign an ID")
	}
	if repo.callCount != 1 {
		t.Errorf("repository was called %d times, want exactly 1", repo.callCount)
	}
}

func TestMovieService_GetByID_RejectsEmptyID(t *testing.T) {
	repo := newFakeRepository()
	svc := NewMovieService(repo)

	_, err := svc.GetByID(context.Background(), "")
	if !errors.Is(err, domain.ErrValidation) {
		t.Fatalf("got err %v, want domain.ErrValidation", err)
	}
}

func TestMovieService_GetByID_PropagatesNotFound(t *testing.T) {
	repo := newFakeRepository()
	svc := NewMovieService(repo)

	_, err := svc.GetByID(context.Background(), "missing")
	if !errors.Is(err, domain.ErrNotFound) {
		t.Fatalf("got err %v, want domain.ErrNotFound", err)
	}
}

func TestMovieService_Create_PropagatesRepositoryFailure(t *testing.T) {
	repo := newFakeRepository()
	repo.forceErr = errors.New("disk on fire")
	svc := NewMovieService(repo)

	
	_, err := svc.GetAll(context.Background())
	if err == nil {
		t.Fatal("expected error to propagate from repository, got nil")
	}
}