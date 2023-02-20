package main

import (
	"net/http"

	"github.com/Gideon-isa/ravefoods/internal/config"
	"github.com/Gideon-isa/ravefoods/internal/handlers"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/v5/middleware"
)

func routes(app *config.AppConfig) http.Handler {
	mux := chi.NewRouter()

	// appends middleware handler to the mux middleware stack
	// This is used to process a request and perform an action before getting to the handlers
	mux.Use(middleware.Recoverer)

	// CSRF check
	mux.Use(NoSurf)

	// saving the session
	mux.Use(SessionLoad)

	// return an handler that contains tye files in the directory
	handler := http.FileServer(http.Dir("./static/"))

	//
	mux.Handle("/static/*", http.StripPrefix("/static", handler))

	mux.Get("/", handlers.Repo.HomePage)
	mux.Get("/signup", handlers.Repo.SignUp)
	mux.Post("/PostSignUp", handlers.Repo.PostSignUp)
	mux.Get("/nigerian-food", handlers.Repo.NigerianFoods)

	// http.HandleFunc("/", handlers.Repo.HomePage)
	// http.HandleFunc("/signup", handlers.Repo.SignUp)
	// http.HandleFunc("/nigerian-food", handlers.Repo.NigerianFoods)

	return mux

}
