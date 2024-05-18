package middlewares

import (
	"io"
	"net/http"
	"context"

	"github.com/gomodule/redigo/redis"

	"github.com/FelineJTD/secure-chat-kripto/server/logger"
	"github.com/FelineJTD/secure-chat-kripto/server/handlers"
)

type contextKey string

const (
	BODY contextKey = "body"
)

func DecryptMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var err error = nil
		defer logger.HandleError(err)

		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		// Decrypt body
		SharedKey, err := handlers.GetSharedKey(r.RemoteAddr)

		if err != nil {
			if err == redis.ErrNil {
				http.Error(w, "No Shared Key Recorded, Must Perform Handshake", http.StatusBadRequest)
				return
			}

			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		plaintext, err := handlers.Decrypt(SharedKey, body)

		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		ctx := context.WithValue(r.Context(), BODY, plaintext)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}