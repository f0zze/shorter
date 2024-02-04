package storage

import "fmt"

func NewInMemoryStorage() (Storage, error) {
	return &URLStorage{
		data: make(map[string]*ShortURL),
	}, nil
}

func (s *URLStorage) Find(uuid string) (*ShortURL, bool) {
	value, ok := s.data[uuid]
	return value, ok
}

func (s *URLStorage) Save(url []ShortURL, _ bool) error {
	for _, u := range url {
		s.data[u.UUID] = &u
	}

	return nil
}

func (s *URLStorage) Ping() bool {
	return true
}

func (s *URLStorage) Size() int {
	return len(s.data)
}

func (s *URLStorage) Close() error {
	return nil
}

func (s *URLStorage) FindShortURLBy(originalURL string) (string, error) {
	// TODO implement logic
	fmt.Print(originalURL)
	return "", nil
}
