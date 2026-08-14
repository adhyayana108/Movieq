package domain

import "errors"

var (
	ErrNotFound   = errors.New("movie not found")
	ErrValidation = errors.New("validation failed")
)

type Director struct {
	FirstName string `json:"firstname"`
	LastName  string `json:"lastname"`
}

type Movie struct {
	ID       string    `json:"id"`
	ISBN     string    `json:"isbn"`
	Title    string    `json:"title"`
	Director *Director `json:"director,omitempty"`
}
