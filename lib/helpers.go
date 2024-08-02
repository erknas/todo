package lib

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

const Layout = "20060102"

type EmptyJSON struct{}

type ApiErr struct {
	Error string `json:"error"`
}

type ApiFunc func(w http.ResponseWriter, r *http.Request) error

func MakeHTTP(fn ApiFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := fn(w, r); err != nil {
			WriteJSON(w, http.StatusBadRequest, ApiErr{Error: err.Error()})
		}
	}
}

func WriteJSON(w http.ResponseWriter, status int, v any) error {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(status)

	return json.NewEncoder(w).Encode(v)
}

func IsDate(s string) bool {
	_, err := time.Parse("02.01.2006", s)
	return err == nil
}

func ParseTime(s string) (string, error) {
	date, err := time.Parse("02.01.2006", s)
	if err != nil {
		return "", fmt.Errorf("invalid date")
	}

	return date.Format(Layout), nil
}
