package services

import (
	"github.com/f0zze/shorter/internal/app/storage"
	"math/rand"
	"time"
)

type ShortURLService struct {
	ResultURL string
	Storage   storage.URLStorage
}

func (service *ShortURLService) CreateNewShortURL(url string) string {
	urlID := generateRandomString(5)

	service.Storage.Set(urlID, url)

	return service.ResultURL + "/" + urlID
}

func (service *ShortURLService) FindURLByID(shortURLID string) string {
	url := service.Storage.Find(shortURLID)

	return url
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
