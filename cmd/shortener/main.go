package main

import (
	"fmt"
	"github.com/f0zze/shorter/internal/app/handlers"
	"net/http"
)

func mainHandler(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodPost:
		handlers.PostHandler(w, req)

	case http.MethodGet:
		handlers.GetHandler(w, req)
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
