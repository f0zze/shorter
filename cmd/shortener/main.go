package main

import (
	"flag"
	"github.com/f0zze/shorter/internal/app/handlers"
	"github.com/f0zze/shorter/internal/app/services"
	"github.com/f0zze/shorter/internal/app/storage"
	chi2 "github.com/go-chi/chi/v5"
	"net/http"
)

type ServerConfig struct {
	host     string
	response string
}

func getConfig() ServerConfig {
	host := flag.String("a", "localhost:8080", "Server URL")
	destHost := flag.String("b", "http://localhost:8080", "Response server URL")
	flag.Parse()

	return ServerConfig{
		*host,
		*destHost,
	}
}

func main() {
	config := getConfig()
	runServer(config)
}

func runServer(config ServerConfig) {
	var urlStorage = storage.NewStorage()
	var shortURLServices = services.ShortURLService{
		ResultURL: config.response,
		Storage:   urlStorage,
	}
	var rootHandler = handlers.RootHandler{
		URLService: shortURLServices,
	}

	router := chi2.NewRouter()

	router.Get("/{id}", rootHandler.GetHandler)
	router.Post("/", rootHandler.PostHandler)

	err := http.ListenAndServe(config.host, router)

	if err != nil {
		panic(err)
	}
}
