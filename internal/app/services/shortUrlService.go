package services

import (
	"github.com/f0zze/shorter/internal/app/storage"
	"math/rand"
	"time"
)

func CreateNewShortURL(url string) string {
	urlID := generateRandomString(5)

	storage.Set(urlID, url)

	return `http://localhost:8080/` + urlID
}

func FindURLByID(shortURLID string) string {
	url := storage.Find(shortURLID)

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
