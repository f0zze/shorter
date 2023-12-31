package handlers

import (
	"github.com/f0zze/shorter/internal/app/services"
	"github.com/go-chi/chi/v5"
	"io"
	"net/http"
)

type RootHandler struct {
	URLService services.ShortURLService
}

func (rootHandler *RootHandler) PostHandler(resp http.ResponseWriter, req *http.Request) {
	body, err := io.ReadAll(req.Body)

	defer req.Body.Close()

	if err != nil {
		http.Error(resp, "Could not read url", http.StatusConflict)
		return
	}

	url := string(body)

	if url == "" {
		http.Error(resp, "Invalid url", http.StatusNotFound)
		return
	}

	shortURL := rootHandler.URLService.CreateNewShortURL(url)

	resp.Header().Add("Content-Type", "text/plain")
	resp.WriteHeader(http.StatusCreated)
	resp.Write([]byte(shortURL))

}

func (rootHandler *RootHandler) GetHandler(resp http.ResponseWriter, req *http.Request) {
	urlID := chi.URLParam(req, "id")

	if urlID == "" {
		http.NotFound(resp, req)
		return
	}

	url := rootHandler.URLService.FindURLByID(urlID)

	if url == "" {
		url = `http://localhost:8080`
	}

	resp.Header().Add("Content-Type", "text/plain")
	http.Redirect(resp, req, url, http.StatusTemporaryRedirect)
}
