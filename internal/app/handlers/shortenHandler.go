package handlers

import (
	"encoding/json"
	"errors"
	"github.com/f0zze/shorter/internal/app/models"
	"github.com/f0zze/shorter/internal/app/services"
	"github.com/f0zze/shorter/internal/app/storage"
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
	var urls []models.OriginalURL

	decoder := json.NewDecoder(req.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&urls); err != nil {
		resp.WriteHeader(http.StatusBadRequest)
		return
	}

	result, err := h.URLService.CreateBatch(urls)

	if err != nil {
		resp.WriteHeader(http.StatusInternalServerError)
		return
	}

	json, err := json.Marshal(result)

	if err != nil {
		resp.WriteHeader(http.StatusInternalServerError)
		return
	}

	resp.Header().Add("Content-Type", "application/json")
	resp.WriteHeader(http.StatusCreated)
	resp.Write(json)
}

func (h *ShortenHandler) Post(resp http.ResponseWriter, req *http.Request) {
	fullURL := FullURL{}

	decoder := json.NewDecoder(req.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&fullURL); err != nil {
		http.Error(resp, err.Error(), http.StatusBadRequest)
		return
	}

	// TODO: handle other errors errors
	shortURL, err := h.URLService.Create(fullURL.URL)

	status := http.StatusCreated
	if errors.Is(err, storage.ErrConflict) {
		status = http.StatusConflict
	}

	responseData := ShortenURL{shortURL}

	shorter, err := json.Marshal(responseData)
	if err != nil {
		http.Error(resp, err.Error(), http.StatusInternalServerError)
		return
	}

	resp.Header().Add("Content-Type", "application/json")
	resp.WriteHeader(status)
	resp.Write(shorter)
}
