package healthchek

import (
	"context"
	"encoding/json"
	"net/http"
)

func liveness(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	status := struct {
		Status string `json:"status"`
	}{
		Status: "ok",
	}
	return json.NewEncoder(w).Encode(status)
}

func readiness(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	status := struct {
		Status string `json:"status"`
	}{
		Status: "ok",
	}
	return json.NewEncoder(w).Encode(status)
}
