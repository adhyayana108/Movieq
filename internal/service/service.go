package service

import(

	"context"
 
	"github.com/adhyayana108/movieq/internal/domain"

)

type MovieService interface{

	GetAll(ctx context.Context) ([]domain.Movie , error)

	GetByID(ctx context.Context , id string) (domain.Movie , error)

	Create(ctx context.Context ,movie domain.Movie) (domain.Movie , error)

	Update(ctx context.Context ,id string , movie domain.Movie) (domain.Movie , error)

	Delete(ctx context.Context , id string) error

}