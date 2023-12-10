package main

import (
	"fmt"
	"github.com/f0zze/shorter/internal/app/handlers"
	"github.com/f0zze/shorter/internal/app/services"
	"github.com/f0zze/shorter/internal/app/storage"
	chi2 "github.com/go-chi/chi/v5"
	"net/http"
)

func main() {
	var urlStorage = storage.NewStorage()
	var shortURLServices = services.ShortURLService{
		Storage: urlStorage,
	}
	var rootHandler = handlers.RootHandler{
		URLService: shortURLServices,
	}

	router := chi2.NewRouter()

	router.Get("/{id}", rootHandler.GetHandler)
	router.Post("/", rootHandler.PostHandler)

	err := http.ListenAndServe(`:8080`, router)

	if err != nil {
		panic(err)
	}

	fmt.Println("Server started...")
}
