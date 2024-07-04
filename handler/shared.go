package handler

import (
	"encoding/json"
	"github.com/a-h/templ"
	"log/slog"
	"net/http"
)

func Make(h func(w http.ResponseWriter, r *http.Request) error) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := h(w, r); err != nil {
			slog.Error("internal server error", "err", err, "path", r.URL.Path)
		}
	}
}

func RespondWithError(w http.ResponseWriter, resError error) {
	responseData := map[string]interface{}{
		"status":  http.StatusInternalServerError,
		"error":   resError.Error(),
		"message": "internal server error",
	}
	responseJSON, err := json.Marshal(responseData)
	if err != nil {
		slog.Error("unable to respond with error", "err", err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusInternalServerError)
	_, err = w.Write(responseJSON)
	if err != nil {
		slog.Error("unable to write response", "err", err)
	}
}

func Respond(w http.ResponseWriter, responseData map[string]interface{}, statusCode int) {
	responseJSON, err := json.Marshal(responseData)
	if err != nil {
		RespondWithError(w, err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	_, err = w.Write(responseJSON)
	if err != nil {
		slog.Error("unable to write response", "err", err)
	}
}

func render(w http.ResponseWriter, r *http.Request, component templ.Component) error {
	return component.Render(r.Context(), w)
}
