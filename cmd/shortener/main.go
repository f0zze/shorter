package main

import (
	"fmt"
	"github.com/f0zze/shorter/internal/app/handlers"
	"github.com/f0zze/shorter/internal/app/services"
	"github.com/f0zze/shorter/internal/app/storage"
	"net/http"
)

var urlStorage = storage.NewStorage()
var shortURLServices = services.ShortURLService{
	Storage: urlStorage,
}
var rootHandler = handlers.RootHandler{
	URLService: shortURLServices,
}

func mainHandler(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodPost:
		rootHandler.PostHandler(w, req)

	case http.MethodGet:
		rootHandler.GetHandler(w, req)
	}
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc(`/`, mainHandler)
	err := http.ListenAndServe(`:8080`, mux)

	if err != nil {
		panic(err)
	}

	fmt.Println("Server started...")
}
