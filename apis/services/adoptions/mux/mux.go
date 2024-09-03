package mux

import (
	"encoding/json"
	"net/http"
)

func WebAPI() *http.ServeMux {

	mux := http.NewServeMux()

	h := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		status := struct {
			Status string `json:"status"`
		}{
			Status: "ok",
		}
		json.NewEncoder(w).Encode(status)
	}
	mux.HandleFunc("/health", h)

	return mux
}
