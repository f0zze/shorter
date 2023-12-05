package storage

var storage = map[string]string{}

func Find(key string) string {
	value := storage[key]

	return value
}

func Set(key string, value string) {
	storage[key] = value
}
