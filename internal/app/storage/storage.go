package storage

import (
	"errors"
	"github.com/f0zze/shorter/cmd/cfg"
	"github.com/f0zze/shorter/internal/app/entity"
)

type ShortURL struct {
	UserID        string
	UUID          string `json:"uuid"`
	ShortURL      string `json:"short_Url"`
	OriginalURL   string `json:"original_Url"`
	CorrelationID string
}

type URLStorage struct {
	data map[string]*ShortURL
}

var ErrConflict = errors.New("data conflict")

type Storage interface {
	Find(uuid string) (*ShortURL, bool)
	FindByUserID(id string) ([]entity.Shorter, error)
	Save(url []ShortURL, strict bool) error
	Size() int
	Ping() bool
	Close() error
	FindShortURLBy(originalURL string) (string, error)
}

func NewStorage(config *cfg.ServerConfig) (Storage, error) {
	if config.DSN != "" {
		return NewPostgresStorage(config.DSN)
	}

	if config.LogFilePath != "" {
		return NewFileStorage(config.FileStoragePath)
	}

	return NewInMemoryStorage()
}
