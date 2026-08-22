package main

import (
	"log"
	"net/http"

	"github.com/adhyayana108/movieq/internal/handler"
	"github.com/adhyayana108/movieq/internal/repository"
	"github.com/adhyayana108/movieq/internal/service"
)

func main() {
	
	repo := repository.NewMemoryRepository()

	svc := service.NewMovieService(repo)

	h := handler.NewMovieHandler(svc)

	mux := handler.NewRouter(h)

	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	log.Println("movieq API listening on http://localhost:8080")

	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}



