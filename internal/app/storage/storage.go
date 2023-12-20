package storage

type URLStorage struct {
	data map[string]string
}

func NewStorage() URLStorage {
	return URLStorage{
		data: make(map[string]string),
	}
}

func (s *URLStorage) Find(key string) string {
	value := s.data[key]
	return value
}

func (s *URLStorage) Set(key string, value string) {
	s.data[key] = value
}

func (s *URLStorage) Size() int {
	return len(s.data)
}
