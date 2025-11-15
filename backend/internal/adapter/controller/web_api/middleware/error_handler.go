package middleware

import (
	"encoding/json"
	"github.com/andreychh/coopera/pkg/errors"
	"net/http"

	"github.com/andreychh/coopera/internal/adapter/controller/web_api/error_mapper"
)

type HandlerFuncWithError func(w http.ResponseWriter, r *http.Request) error

func ErrorHandler(next HandlerFuncWithError) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := next(w, r); err != nil {
			if ve, ok := err.(*errors.ValidationError); ok {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusBadRequest)
				_ = json.NewEncoder(w).Encode(map[string]any{
					"error":   ve.Base.Error(),
					"details": ve.Details,
				})
				return
			}

			code, msg := error_mapper.MapErrorToHTTP(err)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(code)
			_ = json.NewEncoder(w).Encode(map[string]string{"error": msg})
		}
	}
}
