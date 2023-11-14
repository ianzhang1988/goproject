package lambda_service

import (
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	v1 "goproject/prototype/lambda/server/internal/routers/v1"
)

func RegisterHandlers() *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("pong"))
	})

	// api v1
	r.Route("/v1", func(r chi.Router) {
		r.Post("/job", v1.ProbeJob)
	})

	return r
}

func Serve(r *chi.Mux) {
	server := http.Server{
		Addr:              ":3000",
		ReadTimeout:       120 * time.Second,
		ReadHeaderTimeout: 30 * time.Second,
		WriteTimeout:      120 * time.Second,
		IdleTimeout:       300 * time.Second,
	}
	// err := http.ListenAndServe(":3000", r)
	err := server.ListenAndServe()
	if err != nil {
		fmt.Println("serve failed:", err)
	}
}
