package storage

func NewInMemoryStorage() (Storage, error) {
	return &URLStorage{
		data: make(map[string]*ShortURL),
	}, nil
}

func (s *URLStorage) Find(uuid string) (*ShortURL, bool) {
	value, ok := s.data[uuid]
	return value, ok
}

func (s *URLStorage) Save(url *ShortURL) error {
	s.data[url.UUID] = url

	return nil
}

func (s *URLStorage) Size() int {
	return len(s.data)
}
