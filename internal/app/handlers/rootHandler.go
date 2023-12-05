package handlers

import (
	"github.com/f0zze/shorter/internal/app/services"
	"io"
	"net/http"
	"strings"
)

func PostHandler(resp http.ResponseWriter, req *http.Request) {
	body, err := io.ReadAll(req.Body)

	if err != nil {
		http.Error(resp, "Could not read url", http.StatusConflict)
		return
	}

	url := string(body)

	shortURL := services.CreateNewShortURL(url)

	resp.Header().Add("Content-Type", "text/plain")
	resp.WriteHeader(http.StatusCreated)
	resp.Write([]byte(shortURL))

}

func GetHandler(resp http.ResponseWriter, req *http.Request) {
	urlID := parseShorURLID(req.URL.Path)

	if urlID == "" {
		http.NotFound(resp, req)
		return
	}

	url := services.FindURLByID(urlID)

	if url == "" {
		url = `http://localhost:8080`
	}

	resp.Header().Add("Content-Type", "text/plain")
	http.Redirect(resp, req, url, http.StatusTemporaryRedirect)
}

func parseShorURLID(path string) string {
	pathSegments := strings.Split(path, "/")

	if len(pathSegments) != 2 {
		return ""
	}

	return pathSegments[1]
}
