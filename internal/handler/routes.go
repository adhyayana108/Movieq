package handler

import "net/http"

func NewRouter(h *MovieHandler) *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /movies", h.GetMovies)
	mux.HandleFunc("GET /movies/{id}", h.GetMovie)
	mux.HandleFunc("POST /movies", h.CreateMovie)
	mux.HandleFunc("PUT /movies/{id}", h.UpdateMovie)
	mux.HandleFunc("DELETE /movies/{id}", h.DeleteMovie)

	mux.HandleFunc("GET /healthz", healthz)

	return mux
}

func healthz(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(`{"status":"ok"}`))
}
