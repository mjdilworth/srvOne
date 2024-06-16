package server

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
)

func (s server) Health(w http.ResponseWriter, r *http.Request) {

	type response struct {
		Message string `json:"message"`
	}

	if r.Header.Get("content-type") != "" && r.Header.Get("content-type") != "application/json" {
		w.WriteHeader(http.StatusNotImplemented)
		fmt.Printf("content-type not implemented: %s", r.Header.Get("content-type"))
		return
	}

	switch r.Method {
	case http.MethodGet:
		w.Header().Set("content-type", "application/json")
		err := json.NewEncoder(w).Encode(response{Message: "I am healthy"})
		if err != nil {
			// w.WriteHeader(http.StatusInternalServerError)
			slog.Info("error encoding json: %s", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		slog.Info("method not allowed: %s", r.Method)
	}
}
