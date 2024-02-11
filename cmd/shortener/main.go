package main

import (
	"log"
	"net/http"

	chi2 "github.com/go-chi/chi/v5"

	"github.com/f0zze/shorter/cmd/cfg"
	"github.com/f0zze/shorter/internal/app/handlers"
	"github.com/f0zze/shorter/internal/app/logger"
	"github.com/f0zze/shorter/internal/app/middleware"
	"github.com/f0zze/shorter/internal/app/services"
	"github.com/f0zze/shorter/internal/app/storage"
)

func main() {
	config := cfg.GetConfig()
	runServer(config)
}

func runServer(config cfg.ServerConfig) {
	l := logger.NewLogger(config.LogFilePath)
	withLogging := middleware.WithLogging(&l)

	var urlStorage, err = storage.NewStorage(&config)

	if err != nil {
		log.Fatal(err.Error())
	}

	defer urlStorage.Close()

	var shortURLServices = services.ShortURLService{
		ResultURL: config.Response,
		Storage:   urlStorage,
	}
	var rootHandler = handlers.RootHandler{
		URLService: shortURLServices,
	}

	var shorten = handlers.ShortenHandler{
		URLService: shortURLServices,
	}

	var pingHandler = handlers.PingHandler{
		Storage: urlStorage,
	}

	var user = handlers.UserHandler{
		Storage: urlStorage,
		Service: shortURLServices,
	}

	router := chi2.NewRouter().
		With(middleware.WithAuth()).
		With(middleware.GzipMiddleware())

	router.Post("/", withLogging(rootHandler.PostHandler))
	router.Get("/{id}", withLogging(rootHandler.GetHandler))
	router.Post("/api/shorten", withLogging(shorten.Post))
	router.Post("/api/shorten/batch", withLogging(shorten.Batch))
	router.Get("/api/user/urls", user.Urls)
	router.Delete("/api/user/urls", user.DeleteURL)
	router.Get("/ping", pingHandler.Get)

	error := http.ListenAndServe(config.Host, router)
	l.Info().Msg("Server started")
	if error != nil {
		l.Fatal().Err(err).Msg("Server failed to start")
	}
}
