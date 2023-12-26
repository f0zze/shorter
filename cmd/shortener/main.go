package main

import (
	"github.com/f0zze/shorter/internal/app/logger"
	"net/http"

	chi2 "github.com/go-chi/chi/v5"

	"github.com/f0zze/shorter/cmd/cfg"
	"github.com/f0zze/shorter/internal/app/handlers"
	"github.com/f0zze/shorter/internal/app/services"
	"github.com/f0zze/shorter/internal/app/storage"
)

func main() {
	config := cfg.GetConfig()
	runServer(config)
}

func runServer(config cfg.ServerConfig) {
	l := logger.NewLogger(config.LogFilePath)
	withLogging := logger.WithLogging(&l)

	var urlStorage = storage.NewStorage()
	var shortURLServices = services.ShortURLService{
		ResultURL: config.Response,
		Storage:   urlStorage,
	}
	var rootHandler = handlers.RootHandler{
		URLService: shortURLServices,
	}

	router := chi2.NewRouter()

	router.Get("/{id}", withLogging(rootHandler.GetHandler))
	router.Post("/", withLogging(rootHandler.PostHandler))

	err := http.ListenAndServe(config.Host, router)
	l.Info().Msg("Server started")
	if err != nil {
		l.Fatal().Err(err).Msg("Server failed to start")
	}
}
