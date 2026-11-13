package main

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"

	"github.com/yashcherba/pet-api/internal/order"
)

func NewRouter(st order.Store) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /healthz", handleHealthz)
	mux.HandleFunc("GET /v1/orders/{id}", handleGetOrder(st))
	return mux
}

func handleHealthz(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if _, err := io.WriteString(w, `{"status":"ok"}`); err != nil {
		log.Printf("write healthz response: %v", err)
	}
}

func handleGetOrder(st order.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")

		o, err := st.Get(r.Context(), id)
		if errors.Is(err, order.ErrNotFound) {
			writeJSON(w, http.StatusNotFound, map[string]string{"error": "order not found"})
			return
		}
		if err != nil {
			// #nosec G706 -- id is quoted with %q: newlines are escaped, log forging is not possible.
			log.Printf("get order %q: %v", id, err)
			writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "internal"})
			return
		}

		writeJSON(w, http.StatusOK, o)
	}
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(v); err != nil {
		log.Printf("encode response: %v", err)
	}
}
