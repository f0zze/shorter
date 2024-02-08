package middleware

import (
	"context"
	"net/http"
	"time"

	"github.com/f0zze/shorter/internal/app"
	"github.com/f0zze/shorter/internal/app/services"
)

func WithAuth() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			tokenString, err := r.Cookie("Authorization")

			if err != nil {
				newUserID := services.NewUUID()
				token, err := services.BuildJWTString(newUserID)

				if err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					return
				}

				cookie := http.Cookie{
					Name:     "Authorization",
					Value:    token,
					Expires:  time.Now().Add(24 * time.Hour),
					HttpOnly: true,
					Path:     "/",
				}

				http.SetCookie(w, &cookie)
				ctx := context.WithValue(r.Context(), app.UserIDContext, newUserID)
				next.ServeHTTP(w, r.WithContext(ctx))

				return
			}

			userID, err := services.GetUserID(tokenString.Value)

			if err != nil {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), app.UserIDContext, userID)
			next.ServeHTTP(w, r.WithContext(ctx))
		}

		return http.HandlerFunc(fn)
	}
}
