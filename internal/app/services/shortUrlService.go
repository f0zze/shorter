package services

import (
	"math/rand"
	"time"

	"github.com/f0zze/shorter/internal/app/models"
	"github.com/f0zze/shorter/internal/app/storage"
)

type ShortURLService struct {
	ResultURL string
	Storage   storage.Storage
}

func (service *ShortURLService) CreateNewShortURL(originalURL string) string {
	urlID := generateRandomString(5)

	err := service.Storage.Save([]storage.ShortURL{storage.ShortURL{
		UUID:        urlID,
		ShortURL:    urlID,
		OriginalURL: originalURL,
	}})

	if err != nil {
		return ""
	}

	return service.ResultURL + "/" + urlID
}

func (service *ShortURLService) FindOriginalURLByID(uuid string) (*storage.ShortURL, bool) {
	url, ok := service.Storage.Find(uuid)

	return url, ok
}

func (s *ShortURLService) Save(urlsDTO []models.BatchURL) error {
	var urls []storage.ShortURL

	for _, d := range urlsDTO {
		urls = append(urls, storage.ShortURL{
			UUID:        generateRandomString(5),
			ShortURL:    d.CorrelationId,
			OriginalURL: d.OriginalUrl,
		})
	}

	return s.Storage.Save(urls)
}

func generateRandomString(length int) string {
	// Define the characters allowed in a URL
	const urlSafeCharset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

	// Seed the random number generator
	source := rand.NewSource(time.Now().UnixNano())
	rng := rand.New(source)
	// Create a byte slice to hold the random string
	randomString := make([]byte, length)

	// Fill the byte slice with random characters from the URL-safe charset
	for i := range randomString {
		randomString[i] = urlSafeCharset[rng.Intn(len(urlSafeCharset))]
	}

	return string(randomString)

}
