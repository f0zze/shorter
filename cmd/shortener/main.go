package main

import (
	"net/http"

	chi2 "github.com/go-chi/chi/v5"

	cfg "github.com/f0zze/shorter/cmd/cfg"
	"github.com/f0zze/shorter/internal/app/handlers"
	"github.com/f0zze/shorter/internal/app/services"
	"github.com/f0zze/shorter/internal/app/storage"
)

func main() {
	config := cfg.GetConfig()
	runServer(config)
}

func runServer(config cfg.ServerConfig) {
	var urlStorage = storage.NewStorage()
	var shortURLServices = services.ShortURLService{
		ResultURL: config.Response,
		Storage:   urlStorage,
	}
	var rootHandler = handlers.RootHandler{
		URLService: shortURLServices,
	}

	router := chi2.NewRouter()

	router.Get("/{id}", rootHandler.GetHandler)
	router.Post("/", rootHandler.PostHandler)

	err := http.ListenAndServe(config.Host, router)

	if err != nil {
		panic(err)
	}
}
