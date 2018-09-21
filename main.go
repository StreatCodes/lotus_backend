package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"

	_ "github.com/lib/pq"
)

func main() {
	server := GetConfig()

	buildPages(server)

	//Http server setup
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))

	r.Get("/*", ServePageHandler(server))

	r.Route("/api", func(r chi.Router) {
		r.Post("/authorize", authorizeHandler(server))

		r.Post("/users", CreateUserHandler(server))
		r.Post("/componets", CreateComponentHandler(server))
		r.Post("/pages", CreatePageHandler(server))
	})

	fmt.Println("Starting http server on: " + server.HTTPAddress())
	http.ListenAndServe(server.HTTPAddress(), r)
}
