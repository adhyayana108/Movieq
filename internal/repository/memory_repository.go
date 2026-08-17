package repository

import (
	"context"
	"fmt"
	"sort"
	"sync"

	"github.com/adhyayana108/movieq/internal/domain"
)

type MemoryRepository struct {
	mu     sync.RWMutex
	movies map[string]domain.Movie
	nextID int64
}

func NewMemoryRepository() *MemoryRepository {
	return &MemoryRepository{
		movies: make(map[string]domain.Movie),
		nextID: 1,
	}
}

func (r *MemoryRepository) Seed(movies ...domain.Movie) {
	r.mu.Lock()
	defer r.mu.Unlock()
	for _, m := range movies {
		if m.ID == "" {
			m.ID = fmt.Sprintf("%d", r.nextID)
			r.nextID++
		}
		r.movies[m.ID] = m
	}
}

func (r *MemoryRepository) GetAll(_ context.Context) ([]domain.Movie, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	result := make([]domain.Movie, 0, len(r.movies))
	for _, m := range r.movies {
		result = append(result, m)
	}


	sort.Slice(result, func(i, j int) bool { return result[i].ID < result[j].ID })
	return result, nil
}

func (r *MemoryRepository) GetByID(_ context.Context, id string) (domain.Movie, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	m, ok := r.movies[id]
	if !ok {
		return domain.Movie{}, domain.ErrNotFound
	}
	return m, nil
}

func (r *MemoryRepository) Create(_ context.Context, movie domain.Movie) (domain.Movie, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	movie.ID = fmt.Sprintf("%d", r.nextID)
	r.nextID++
	r.movies[movie.ID] = movie
	return movie, nil
	
}

func (r *MemoryRepository) Update(_ context.Context, id string, movie domain.Movie) (domain.Movie, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.movies[id]; !ok {
		return domain.Movie{}, domain.ErrNotFound
	}
	movie.ID = id
	r.movies[id] = movie
	return movie, nil
}

func (r *MemoryRepository) Delete(_ context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.movies[id]; !ok {
		return domain.ErrNotFound
	}
	delete(r.movies, id)
	return nil
}


var _ MovieRepository = (*MemoryRepository)(nil)

