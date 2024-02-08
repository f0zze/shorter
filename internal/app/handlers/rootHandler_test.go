package handlers

import (
	"context"
	"github.com/f0zze/shorter/internal/app"
	"github.com/f0zze/shorter/internal/app/storage"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestPostHandler(t *testing.T) {
	urlStorage, err := storage.NewInMemoryStorage()
	if err != nil {
		panic("Could not create in memory storage")
	}
	handlers := RootHandler{
		URLService: struct {
			ResultURL string
			Storage   storage.Storage
		}{Storage: urlStorage, ResultURL: "http://localhost:8888"},
	}

	t.Run("should return new shorter", func(t *testing.T) {
		// Create a mock user ID
		mockUserID := "mockUserID"

		// Create a context with the mock user ID
		ctx := context.WithValue(context.Background(), app.UserIDContext, mockUserID)

		// Create a request with the mock context
		request := httptest.NewRequest(http.MethodPost, "/", strings.NewReader("https://yandex.ru"))
		request = request.WithContext(ctx)
		record := httptest.NewRecorder()

		handlers.PostHandler(record, request)

		response := record.Result()
		body, err := io.ReadAll(response.Body)

		if err != nil {
			panic("Could not parse response body")
		}
		defer response.Body.Close()

		assert.Contains(t, string(body), "http://localhost:8888/")
		assert.Equal(t, response.StatusCode, http.StatusCreated)
		assert.Equal(t, 1, urlStorage.Size())
	})

	t.Run("should return 404 when send empty body", func(t *testing.T) {
		// Create a mock user ID
		mockUserID := "mockUserID"

		// Create a context with the mock user ID
		ctx := context.WithValue(context.Background(), app.UserIDContext, mockUserID)

		request := httptest.NewRequest(http.MethodPost, "/", nil)
		record := httptest.NewRecorder()
		request.WithContext(ctx)

		handlers.PostHandler(record, request)

		response := record.Result()
		defer response.Body.Close()

		assert.Equal(t, response.StatusCode, http.StatusNotFound)
	})
}

func TestGetHandler(t *testing.T) {
	urlStorage, err := storage.NewInMemoryStorage()

	if err != nil {
		panic("Could not create in memory storage")
	}

	handlers := RootHandler{
		URLService: struct {
			ResultURL string
			Storage   storage.Storage
		}{Storage: urlStorage, ResultURL: "http://localhost:2222"},
	}
	t.Run("should redirect to full url by requested id", func(t *testing.T) {
		urlToSave := "https://yandex.ru"
		url, _ := handlers.URLService.CreateURL(urlToSave, "testid")
		urlID := strings.Split(url, "/")[3]

		request := httptest.NewRequest(http.MethodGet, "/{id}", nil)
		record := httptest.NewRecorder()

		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("id", urlID)

		request = request.WithContext(context.WithValue(request.Context(), chi.RouteCtxKey, rctx))

		handlers.GetHandler(record, request)

		response := record.Result()
		defer response.Body.Close()

		assert.Equal(t, http.StatusTemporaryRedirect, response.StatusCode)
		assert.Equal(t, response.Header.Get("Location"), urlToSave)
	})

	t.Run("should return full url with requested id ", func(t *testing.T) {
		request := httptest.NewRequest(http.MethodGet, "/", nil)
		record := httptest.NewRecorder()

		handlers.GetHandler(record, request)

		response := record.Result()
		defer response.Body.Close()

		assert.Equal(t, response.StatusCode, http.StatusNotFound)
	})

}
