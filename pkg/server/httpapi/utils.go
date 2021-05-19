package httpapi

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

// respondJSON writes data as JSON
func respondJSON(w http.ResponseWriter, data interface{}, statusCode int) error {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(statusCode)

	if data == nil {
		return errors.New("can't respond with nil data")
	}

	if err := json.NewEncoder(w).Encode(data); err != nil {
		return fmt.Errorf("json encode: %w", err)
	}

	return nil
}

// respondJSONError responds with error message and specified HTTP status
func respondJSONError(w http.ResponseWriter, err error, statusCode int) error {
	if err == nil {
		return errors.New("can't respond with nil err")
	}
	return respondJSON(w,
		struct {
			Error string `json:"error"`
		}{Error: err.Error()},
		statusCode)
}

// decodeJSON decodes data from request body as JSON
func decodeJSON(r *http.Request, v interface{}) error {
	return json.NewDecoder(r.Body).Decode(v)
}
