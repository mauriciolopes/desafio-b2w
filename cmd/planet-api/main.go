package main

import (
	"desafio-b2w/core"
	"fmt"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/rs/cors"
)

func main() {
	f := core.Init()
	f.Configure("http-addr", ":8080", "Endereço e porta da API REST")
	f.RegisterShutdown(func() {
		fmt.Println("Que a força esteja com você!")
	})
	f.Run(func(f *core.Facade) {
		fmt.Println("Sistema iniciado. Pressione Ctrl+C para encerrar.")
		r := chi.NewRouter()
		r.Use(middleware.RealIP)
		cors := cors.New(cors.Options{
			AllowedOrigins:   []string{"*"},
			AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
			AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token", "X-Correlation-ID"},
			ExposedHeaders:   []string{"Link"},
			AllowCredentials: false,
			MaxAge:           300,
		})
		r.Use(cors.Handler)
		h := &http.Server{Addr: f.Get("http-addr"), Handler: rotas(f)}
		h.ListenAndServe()
	})
}
