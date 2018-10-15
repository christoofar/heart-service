package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"

	"heart-service/floors"
	"heart-service/heart"
	"heart-service/ping"
	"heart-service/steps"
)

func Routes() *chi.Mux {
	router := chi.NewRouter()
	router.Use(
		render.SetContentType(render.ContentTypeJSON),
		middleware.Logger,
		middleware.DefaultCompress,
		middleware.RedirectSlashes,
		middleware.Recoverer,
	)

	router.Route("/v1/", func(r chi.Router) {
		r.Mount("/api/ping", ping.Routes())
		r.Mount("/api/heart", heart.Routes())
		r.Mount("/api/steps", steps.Routes())
		r.Mount("/api/floors", floors.Routes())
	})

	return router
}

func main() {
	fmt.Println("Heart Monitor Cache Service - A service that stores biometric data for pickup.")
	log.Println("[Heartmon] Heart Monitor Cache Service - A service that stores biometric data for pickup.")

	router := Routes()

	walkFunc := func(method string, route string, handler http.Handler, middlewares ...func(http.Handler) http.Handler) error {
		log.Printf("%s %s\n", method, route) // Walk and print out all routes
		return nil
	}

	if err := chi.Walk(router, walkFunc); err != nil {
		log.Panicf("Logging err: %s\n", err.Error()) // panic if there is an error
	}

	log.Fatal(http.ListenAndServe(":8080", router))

}
