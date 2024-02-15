package handlers

import (
	"errors"
	"github.com/f0zze/shorter/internal/app"
	"io"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/f0zze/shorter/internal/app/services"
	"github.com/f0zze/shorter/internal/app/storage"
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

	userID := req.Context().Value(app.UserIDContext).(string)
	shortURL, err := rootHandler.URLService.CreateURL(url, userID)

	status := http.StatusCreated
	if errors.Is(err, storage.ErrConflict) {
		status = http.StatusConflict
	} else if err != nil {
		http.Error(resp, "Create url failed", http.StatusInternalServerError)
		return
	}

	resp.Header().Add("Content-Type", "text/plain")
	resp.WriteHeader(status)
	resp.Write([]byte(shortURL))

}

func (rootHandler *RootHandler) GetHandler(resp http.ResponseWriter, req *http.Request) {
	urlID := chi.URLParam(req, "id")

	if urlID == "" {
		http.NotFound(resp, req)
		return
	}

	url, err := rootHandler.URLService.FindURL(urlID)

	if errors.Is(services.URLDeletedErr, err) {
		http.Error(resp, "Deleted", http.StatusGone)
		return
	}

	redirectURL := `http://localhost:8080`

	if err == nil {
		redirectURL = url.OriginalURL
	}

	resp.Header().Add("Content-Type", "text/plain")
	http.Redirect(resp, req, redirectURL, http.StatusTemporaryRedirect)
}
