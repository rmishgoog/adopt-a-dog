package healthchek

import (
	"encoding/json"
	"net/http"
)

func liveness(w http.ResponseWriter, r *http.Request) {
	status := struct {
		Status string `json:"status"`
	}{
		Status: "ok",
	}
	json.NewEncoder(w).Encode(status)
}

func readiness(w http.ResponseWriter, r *http.Request) {
	status := struct {
		Status string `json:"status"`
	}{
		Status: "ok",
	}
	json.NewEncoder(w).Encode(status)
}
