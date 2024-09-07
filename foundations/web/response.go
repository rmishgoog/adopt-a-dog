package web

import (
	"context"
	"fmt"
	"net/http"

	"github.com/go-json-experiment/json"
)

func Respond(ctx context.Context, w http.ResponseWriter, r *http.Request, data any, statusCode int) error {

	if statusCode == http.StatusNoContent {
		w.WriteHeader(statusCode)
		return nil
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("marshaling data: %w", err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	if _, err := w.Write(jsonData); err != nil {
		return fmt.Errorf("writing response: %w", err)
	}
	return nil
}
