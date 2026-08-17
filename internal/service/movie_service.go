package service
 
import (
	"context"
	"fmt"
	"strings"
 
	"github.com/adhyayana108/movieq/internal/domain"
	"github.com/adhyayana108/movieq/internal/repository"
)

type movieService struct{
	repo repository.MovieRepository
}

func NewMovieService(repo repository.MovieRepository) MovieService{
	return &movieService{repo: repo}
}


func(s *movieService) GetAll(ctx context.Context) ([]domain.Movie, error){
	return s.repo.GetAll(ctx)
}


func (s *movieService) GetByID(ctx context.Context, id string) (domain.Movie, error) {
	if strings.TrimSpace(id) == "" {
		return domain.Movie{}, fmt.Errorf("%w: id is required", domain.ErrValidation)
	}
	return s.repo.GetByID(ctx, id)
}


func (s *movieService) Create(ctx context.Context, movie domain.Movie) (domain.Movie, error) {
	if err := validate(movie); err != nil {
		return domain.Movie{}, err
	}
	return s.repo.Create(ctx, movie)
}


func (s *movieService) Update(ctx context.Context, id string, movie domain.Movie) (domain.Movie, error) {
	if strings.TrimSpace(id) == "" {
		return domain.Movie{}, fmt.Errorf("%w: id is required", domain.ErrValidation)
	}
	if err := validate(movie); err != nil {
		return domain.Movie{}, err
	}
	return s.repo.Update(ctx, id, movie)
}


func (s *movieService) Delete(ctx context.Context, id string) error {
	if strings.TrimSpace(id) == "" {
		return fmt.Errorf("%w: id is required", domain.ErrValidation)
	}
	return s.repo.Delete(ctx, id)
}


func validate(m domain.Movie) error {
	var problems []string

	if strings.TrimSpace(m.Title) == "" {
		problems = append(problems, "title is required")
	}
	if strings.TrimSpace(m.ISBN) == "" {
		problems = append(problems, "isbn is required")
	}
	if m.Director != nil {
		if strings.TrimSpace(m.Director.FirstName) == "" && strings.TrimSpace(m.Director.LastName) == "" {
			problems = append(problems, "director, if provided, needs at least one name field")
		}
	}

	if len(problems) > 0 {
		return fmt.Errorf("%w: %s", domain.ErrValidation, strings.Join(problems, " ; "))
	}
	return nil
}


