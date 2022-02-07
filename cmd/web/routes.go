package main

import (
	"net/http"

	"github.com/javalen/bookings/internal/config"
	"github.com/javalen/bookings/internal/handlers"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func routes(app *config.AppConfig) http.Handler {

	mux := chi.NewRouter()

	mux.Use(middleware.Recoverer)
	mux.Use(NoSurf)
	mux.Use(SessionLoad)

	mux.Get("/", handlers.Repo.Home)
	mux.Get("/about", handlers.Repo.About)
	mux.Get("/generals", handlers.Repo.Generals)
	mux.Get("/majors", handlers.Repo.Majors)
	mux.Get("/search", handlers.Repo.Search)
	mux.Get("/reserve", handlers.Repo.Reservation)
	mux.Get("/contact", handlers.Repo.Contact)
	mux.Get("/summary", handlers.Repo.Summary)

	mux.Post("/search", handlers.Repo.PostSearch)
	mux.Post("/search-availability", handlers.Repo.AvailabilityPost)
	mux.Post("/reserve", handlers.Repo.PostReservation)

	fileServer := http.FileServer(http.Dir("./static/"))

	mux.Handle("/static/*", http.StripPrefix("/static", fileServer))
	return mux
}
