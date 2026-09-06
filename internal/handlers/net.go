package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/vonrimez/TaskAPI/domain"
)

// impements basic response structure
// contains: status code, headers (only with "Content-Type: application/json" yet) and body - any data which can be encoded to json)
type Response struct {
	statusCode int
	headers    map[string]string
	body       any
}

func setValue(r *http.Request, key any, value any) {
	*r = *r.WithContext(context.WithValue(r.Context(), key, value))
}

func getValue(r *http.Request, key any) any {
	return r.Context().Value(key)
}

func respondWithJSON(w http.ResponseWriter, resp *Response) error {
	w.Header().Set("Content-Type", "application/json")
	for key, value := range resp.headers {
		w.Header().Set(key, value)
	}
	w.WriteHeader(resp.statusCode)
	err := json.NewEncoder(w).Encode(resp.body)
	return err
}

func respondWithAppError(w http.ResponseWriter, err error) {
	var appErr *domain.AppError
	if errors.As(err, &appErr) {
		respondWithJSON(w, &Response{
			statusCode: appErr.Code,
			body:       H{"error": appErr.Err.Error()},
		})
		return
	}
	respondWithJSON(w, &Response{
		statusCode: http.StatusInternalServerError,
		body:       H{"error": "internal error"},
	})
}
