package storage

import "github.com/f0zze/shorter/cmd/cfg"

type ShortURL struct {
	UUID        string `json:"uuid"`
	ShortURL    string `json:"short_Url"`
	OriginalURL string `json:"original_Url"`
}

type URLStorage struct {
	data map[string]*ShortURL
}

type Storage interface {
	Find(uuid string) (*ShortURL, bool)
	Save(url *ShortURL) error
	Size() int
}

func NewStorage(config *cfg.ServerConfig) (Storage, error) {
	if config.LogFilePath == "" {
		return NewInMemoryStorage()
	}

	return NewFileStorage(config.FileStoragePath)
}
