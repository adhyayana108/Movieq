package repository

import (
	"context"
	"errors"
	"sync"
	"testing"

	"github.com/adhyayana108/movieq/internal/domain"
)

func TestMemoryRepository_CreateAndGetByID(t *testing.T) {
	repo := NewMemoryRepository()
	ctx := context.Background()

	created, err := repo.Create(ctx, domain.Movie{Title: "Arrival", ISBN: "123"})
	if err != nil {
		t.Fatalf("Create returned unexpected error: %v", err)
	}
	if created.ID == "" {
		t.Fatal("Create did not assign an ID")
	}

	got, err := repo.GetByID(ctx, created.ID)
	if err != nil {
		t.Fatalf("GetByID returned unexpected error: %v", err)
	}
	if got.Title != "Arrival" {
		t.Errorf("got title %q, want %q", got.Title, "Arrival")
	}
}

func TestMemoryRepository_GetByID_NotFound(t *testing.T) {
	repo := NewMemoryRepository()

	_, err := repo.GetByID(context.Background(), "does-not-exist")
	if !errors.Is(err, domain.ErrNotFound) {
		t.Errorf("got err %v, want domain.ErrNotFound", err)
	}
}

func TestMemoryRepository_GetAll_NeverReturnsNil(t *testing.T) {
	repo := NewMemoryRepository()

	movies, err := repo.GetAll(context.Background())
	if err != nil {
		t.Fatalf("GetAll returned unexpected error: %v", err)
	}
	if movies == nil {
		t.Error("GetAll returned nil slice; want empty non-nil slice so JSON encodes as [] not null")
	}
	if len(movies) != 0 {
		t.Errorf("got %d movies, want 0", len(movies))
	}
}

func TestMemoryRepository_GetAll_IsSortedByID(t *testing.T) {
	repo := NewMemoryRepository()
	ctx := context.Background()

	for i := 0; i < 5; i++ {
		if _, err := repo.Create(ctx, domain.Movie{Title: "x", ISBN: "y"}); err != nil {
			t.Fatalf("Create returned unexpected error: %v", err)
		}
	}

	movies, err := repo.GetAll(ctx)
	if err != nil {
		t.Fatalf("GetAll returned unexpected error: %v", err)
	}
	for i := 1; i < len(movies); i++ {
		if movies[i-1].ID >= movies[i].ID {
			t.Errorf("movies not sorted: %q came before %q", movies[i-1].ID, movies[i].ID)
		}
	}
}

func TestMemoryRepository_Update(t *testing.T) {
	repo := NewMemoryRepository()
	ctx := context.Background()

	created, _ := repo.Create(ctx, domain.Movie{Title: "Old Title", ISBN: "1"})

	updated, err := repo.Update(ctx, created.ID, domain.Movie{Title: "New Title", ISBN: "1"})
	if err != nil {
		t.Fatalf("Update returned unexpected error: %v", err)
	}
	if updated.Title != "New Title" {
		t.Errorf("got title %q, want %q", updated.Title, "New Title")
	}
	if updated.ID != created.ID {
		t.Errorf("Update changed the ID: got %q, want %q", updated.ID, created.ID)
	}
}

func TestMemoryRepository_Update_NotFound(t *testing.T) {
	repo := NewMemoryRepository()

	_, err := repo.Update(context.Background(), "missing", domain.Movie{Title: "x", ISBN: "y"})
	if !errors.Is(err, domain.ErrNotFound) {
		t.Errorf("got err %v, want domain.ErrNotFound", err)
	}
}

func TestMemoryRepository_Delete(t *testing.T) {
	repo := NewMemoryRepository()
	ctx := context.Background()

	created, _ := repo.Create(ctx, domain.Movie{Title: "x", ISBN: "y"})

	if err := repo.Delete(ctx, created.ID); err != nil {
		t.Fatalf("Delete returned unexpected error: %v", err)
	}
	if _, err := repo.GetByID(ctx, created.ID); !errors.Is(err, domain.ErrNotFound) {
		t.Errorf("movie still retrievable after Delete: err = %v", err)
	}
}

func TestMemoryRepository_Delete_NotFound(t *testing.T) {
	repo := NewMemoryRepository()

	err := repo.Delete(context.Background(), "missing")
	if !errors.Is(err, domain.ErrNotFound) {
		t.Errorf("got err %v, want domain.ErrNotFound", err)
	}
}

func TestMemoryRepository_ConcurrentAccess(t *testing.T) {
	repo := NewMemoryRepository()
	ctx := context.Background()

	var wg sync.WaitGroup
	const workers = 50

	wg.Add(workers)
	for i := 0; i < workers; i++ {
		go func() {
			defer wg.Done()
			m, err := repo.Create(ctx, domain.Movie{Title: "concurrent", ISBN: "z"})
			if err != nil {
				t.Errorf("Create returned unexpected error: %v", err)
				return
			}
			if _, err := repo.GetByID(ctx, m.ID); err != nil {
				t.Errorf("GetByID returned unexpected error: %v", err)
			}
			if _, err := repo.GetAll(ctx); err != nil {
				t.Errorf("GetAll returned unexpected error: %v", err)
			}
		}()
	}
	wg.Wait()

	movies, err := repo.GetAll(ctx)
	if err != nil {
		t.Fatalf("GetAll returned unexpected error: %v", err)
	}
	if len(movies) != workers {
		t.Errorf("got %d movies after concurrent creates, want %d", len(movies), workers)
	}
}
