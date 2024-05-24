package api

import (
	"encoding/json"
	"log/slog"
	"net/http"
)

type Error struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
}

func WriteOk(w http.ResponseWriter, data any) {
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(data); err != nil {
		slog.With("error", err).With("source", "api").Error("failed to encode")
	}
}

func WriteError(w http.ResponseWriter, statusCode int, e *Error) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(e); err != nil {
		slog.With("error", err).With("source", "api").Error("failed to encode")
	}
}
