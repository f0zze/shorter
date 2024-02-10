package middleware

import (
	"context"
	"errors"
	"fmt"
	"github.com/f0zze/shorter/internal/app"
	"github.com/f0zze/shorter/internal/app/services"
	"net/http"
	"time"
)

func createAuthCookie(token string) *http.Cookie {
	return &http.Cookie{
		Name:     "ID",
		Value:    token,
		Expires:  time.Now().Add(24 * time.Hour),
		HttpOnly: true,
		Path:     "/",
	}
}

func setUserIDToContext(r *http.Request, userID string) *http.Request {
	ctx := context.WithValue(r.Context(), app.UserIDContext, userID)
	r.WithContext(ctx)

	return r.WithContext(ctx)
}

func WithAuth() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			tokenString, err := r.Cookie("ID")

			if (errors.Is(http.ErrNoCookie, err)) && r.URL.Path != "/api/user/urls" {
				fmt.Println("Generate new token")
				newUserID := services.NewUUID()
				token, err := services.BuildJWTString(newUserID)

				if err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					return
				}

				http.SetCookie(w, createAuthCookie(token))

				req := setUserIDToContext(r, newUserID)

				next.ServeHTTP(w, req)

				return
			}

			userID := services.GetUserID(tokenString.Value)

			next.ServeHTTP(w, setUserIDToContext(r, userID))
		}

		return http.HandlerFunc(fn)
	}
}
