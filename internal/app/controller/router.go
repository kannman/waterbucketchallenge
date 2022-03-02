package controller

import "github.com/go-chi/chi"

func MountRoutes(r *chi.Mux) {
	r.Post("/v1/search-solution", FindSolutionHandler)
}
