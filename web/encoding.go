package web

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func Encode[T any](w http.ResponseWriter, status int, v T) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(v); err != nil {
		return fmt.Errorf("json encode error: %w", err)
	}

	return nil
}

func Decode[T any](r *http.Request) (T, error) {
	var v T
	if err := json.NewDecoder(r.Body).Decode(&v); err != nil {
		return v, Error{
			Message:    "error decoding json",
			StatusCode: http.StatusBadRequest,
		}
	}

	return v, nil
}

func DecodeValidated[T any](v *Validator, r *http.Request) (T, error) {
	t, err := Decode[T](r)
	if err != nil {
		return t, err
	}

	problems, err := v.Struct(t)
	if err != nil {
		return t, err
	}
	if problems != nil {
		return t, Error{
			Message:    "Validation Error",
			StatusCode: http.StatusBadRequest,
			Data:       problems,
		}
	}

	return t, nil
}
