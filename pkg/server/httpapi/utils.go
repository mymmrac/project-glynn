package httpapi

import (
	"encoding/json"
	"errors"
	"net/http"
)

// respondJSON writes data as JSON
func respondJSON(w http.ResponseWriter, data interface{}, statusCode int) error {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(statusCode)

	if data != nil {
		return json.NewEncoder(w).Encode(data)
	}

	return errors.New("can't respond with nil data")
}

// respondJSONError responds with error message and specified HTTP status
func respondJSONError(w http.ResponseWriter, err error, statusCode int) error {
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
