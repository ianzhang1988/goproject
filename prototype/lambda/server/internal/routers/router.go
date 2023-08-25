package lambda_service

import (
	"fmt"
	"net/http"

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
	
	err := http.ListenAndServe(":3000", r)
	if err != nil {
		fmt.Println("serve failed:", err)
	}
}
