package models

type OriginalURL struct {
	CorrelationID string `json:"correlation_id"`
	OriginalURL   string `json:"original_url"`
}

type ShortURL struct {
	CorrelationID string `json:"correlation_id"`
	ShortURL      string `json:"short_url"`
}

type UserShorter struct {
	OriginalURL string `json:"original_url"`
	ShortURL    string `json:"short_url"`
}
