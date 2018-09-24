package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"

	_ "github.com/lib/pq"
)

func main() {
	server := GetConfig()

	siteTree, err := buildPages(server)
	if err != nil {
		log.Fatalf("Error building site %s", err)
	}
	server.SiteTree = siteTree

	//Http server setup
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))

	r.Get("/*", ServePageHandler(server))

	fs := http.FileServer(http.Dir("./admin"))
	r.Get("/admin/*", http.StripPrefix("/admin", fs).ServeHTTP)

	r.Route("/api", func(r chi.Router) {
		r.Post("/authorize", authorizeHandler(server))

		r.Post("/users", CreateUserHandler(server))
		r.Post("/componets", CreateComponentHandler(server))
		r.Post("/pages", CreatePageHandler(server))
	})

	fmt.Println("Starting http server on: " + server.HTTPAddress())
	http.ListenAndServe(server.HTTPAddress(), r)
}
