package handlers

import (
	"encoding/json"
	"github.com/f0zze/shorter/internal/app/models"
	"github.com/f0zze/shorter/internal/app/services"
	"net/http"
)

type ShortenHandler struct {
	URLService services.ShortURLService
}

type FullURL struct {
	URL string `json:"url"`
}

type ShortenURL struct {
	Result string `json:"result"`
}

func (h *ShortenHandler) Batch(resp http.ResponseWriter, req *http.Request) {
	var urls []models.BatchURL

	decoder := json.NewDecoder(req.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&urls); err != nil {
		resp.WriteHeader(http.StatusBadRequest)
		return
	}

	h.URLService.Save(urls)

	resp.WriteHeader(http.StatusCreated)
}

func (shortenHandler *ShortenHandler) Post(resp http.ResponseWriter, req *http.Request) {
	fullURL := FullURL{}

	decoder := json.NewDecoder(req.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&fullURL); err != nil {
		http.Error(resp, err.Error(), http.StatusBadRequest)
		return
	}

	shortURL := shortenHandler.URLService.CreateNewShortURL(fullURL.URL)
	responseData := ShortenURL{shortURL}

	shorter, err := json.Marshal(responseData)
	if err != nil {
		http.Error(resp, err.Error(), http.StatusInternalServerError)
		return
	}

	resp.Header().Add("Content-Type", "application/json")
	resp.WriteHeader(http.StatusCreated)
	resp.Write(shorter)
}
