package services

import (
	"errors"
	"github.com/jackc/pgx/v5/pgconn"
	"math/rand"
	"time"

	"github.com/f0zze/shorter/internal/app/models"
	"github.com/f0zze/shorter/internal/app/storage"
	"github.com/jackc/pgerrcode"
)

type ShortURLService struct {
	ResultURL string
	Storage   storage.Storage
}

func (s *ShortURLService) CreateURLs(urls []models.OriginalURL) ([]models.ShortURL, error) {

	var data []storage.ShortURL

	for _, u := range urls {
		uuid := NewUUID()
		shortURL := NewShortURL()

		data = append(data, storage.ShortURL{
			UUID:          uuid,
			ShortURL:      shortURL,
			OriginalURL:   u.OriginalURL,
			CorrelationID: u.CorrelationID,
		})
	}

	err := s.Storage.Save(data, false)

	if err != nil {
		return nil, err
	}

	var result []models.ShortURL

	for _, d := range data {
		result = append(result, models.ShortURL{
			ShortURL:      s.ResultURL + "/" + d.ShortURL,
			CorrelationID: d.CorrelationID,
		})
	}

	return result, nil
}

func (s *ShortURLService) CreateURL(originalURL string) (string, error) {
	urlID := NewShortURL()

	err := s.Storage.Save([]storage.ShortURL{storage.ShortURL{
		UUID:        urlID,
		ShortURL:    urlID,
		OriginalURL: originalURL,
	}}, true)

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgerrcode.IsIntegrityConstraintViolation(pgErr.Code) {
			result, err := s.Storage.FindShortURLBy(originalURL)

			if err != nil {
				return "", err
			}

			return s.ResultURL + "/" + result, storage.ErrConflict
		}
	}

	return s.ResultURL + "/" + urlID, err
}

func (s *ShortURLService) FindURL(uuid string) (*storage.ShortURL, bool) {
	url, ok := s.Storage.Find(uuid)

	return url, ok
}

func NewUUID() string {
	return generateRandomString(5)
}

func NewShortURL() string {
	return generateRandomString(5)
}

func generateRandomString(length int) string {
	// Define the characters allowed in a URL
	const urlSafeCharset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

	// Seed the random number generator
	source := rand.NewSource(time.Now().UnixNano())
	rng := rand.New(source)
	// CreateURLs a byte slice to hold the random string
	randomString := make([]byte, length)

	// Fill the byte slice with random characters from the URL-safe charset
	for i := range randomString {
		randomString[i] = urlSafeCharset[rng.Intn(len(urlSafeCharset))]
	}

	return string(randomString)

}
