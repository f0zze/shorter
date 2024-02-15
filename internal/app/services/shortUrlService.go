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

func (s *ShortURLService) CreateURLs(urls []models.OriginalURL, userID string) ([]models.ShortURL, error) {

	var data []storage.ShortURL

	for _, u := range urls {
		uuid := NewUUID()
		shortURL := NewShortURL()

		data = append(data, storage.ShortURL{
			UUID:          uuid,
			ShortURL:      shortURL,
			OriginalURL:   u.OriginalURL,
			CorrelationID: u.CorrelationID,
			UserID:        userID,
		})
	}

	err := s.Storage.Save(data)

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

func (s *ShortURLService) CreateURL(originalURL string, userID string) (string, error) {
	urlID := NewShortURL()

	err := s.Storage.Save([]storage.ShortURL{storage.ShortURL{
		UUID:        urlID,
		ShortURL:    urlID,
		OriginalURL: originalURL,
		UserID:      userID,
	}})

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

var NotFoundErr = errors.New("url not found")
var URLDeletedErr = errors.New("url deleted")

func (s *ShortURLService) FindURL(uuid string) (*storage.ShortURL, error) {
	url, ok := s.Storage.Find(uuid)

	if !ok {
		return nil, NotFoundErr
	}

	if ok && url.DeletedFlag {
		return nil, URLDeletedErr
	}

	return url, nil
}

func (s *ShortURLService) FindByUser(userID string) ([]models.UserShorter, error) {
	urls, err := s.Storage.FindByUserID(userID)

	if err != nil {
		return nil, err
	}

	var result []models.UserShorter

	for _, v := range urls {
		result = append(result, models.UserShorter{OriginalURL: v.OriginalURL, ShortURL: s.addOrigin(v.ShortURL)})
	}

	return result, nil
}

func (s *ShortURLService) DeleteURL(url []string, userID string) error {

	err := s.Storage.DeleteURLsByUserID(url, userID)

	return err
}

func (s *ShortURLService) addOrigin(shortURL string) string {
	return s.ResultURL + "/" + shortURL
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
