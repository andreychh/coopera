package middleware

import (
	"context"
	"net/http"
	"time"
)

func TimeoutMiddleware(duration time.Duration) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx, cancel := context.WithTimeout(r.Context(), duration)
			defer cancel()

			done := make(chan struct{})
			go func() {
				next.ServeHTTP(w, r.WithContext(ctx))
				close(done)
			}()

			select {
			case <-ctx.Done():
				http.Error(w, "request timeout", http.StatusGatewayTimeout)
			case <-done:
			}
		})
	}
}
