package middleware

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/DanielTitkov/netology-slowly/internal/api"
	"github.com/DanielTitkov/netology-slowly/internal/configs"
)

// Middleware holds standart middleware signature.
// Used for brevity.
type Middleware func(http.Handler) http.Handler

// NewTimeout returns timeout middleware
// which stops the requests that are taking too long.
// This is implemented due to the technical requirements,
// though using middleware for the issue seems to be not optimal.
func NewTimeout(cfg configs.Config) Middleware {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx, cancel := context.WithTimeout(
				context.Background(),
				time.Duration(cfg.MaxSlowTimeout)*time.Millisecond,
			)
			defer cancel()

			r = r.WithContext(ctx)
			done := make(chan struct{})
			go func() {
				h.ServeHTTP(w, r)
				close(done)
			}()

			select {
			case <-done:
			case <-ctx.Done():
				resp := api.ErrorResponseBody{
					Error: api.ErrorTooLongTimeout,
				}
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusBadRequest)
				_ = json.NewEncoder(w).Encode(resp)
			}
		})
	}
}
