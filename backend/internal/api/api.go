package api

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/mingassing/app/backend/internal/store"
)

// Handler helps read and write access in api routes
// allows for some cool method shenanigans where we can
// write methods onto the type, and then access something similar
// to self.db in other langs, but essentially keeps our open connections
// to the db explicit and clean. In another file we can define
// `listVehicles` on type handler, and then here we can call `h.listVehicles`
type handler struct {
	db *store.Store
}

// Routes holds the structure of the backend, including the nesting for api
// url routes
func Routes(r chi.Router, db *store.Store) {
	h := &handler{db: db}

	r.Get("/healthz", h.health)

	r.Route("/api/v1", func(r chi.Router) {
		r.Route("/vehicles", func(r chi.Router) {
			r.Post("/", h.createVehicle)
			r.Get("/", h.listVehicles)
			r.Get("/{id}", h.getVehicle)
			r.Delete("/{id}", h.deleteVehicle)
		})
	})
}

// Health check api endpoint, used to determine whether the api is ready to
// receive requests (i.e. the database is reachable)
func (h *handler) health(w http.ResponseWriter, r *http.Request) {
	if err := h.db.Ping(r.Context()); err != nil {
		respondError(w, http.StatusServiceUnavailable, "database unreachable")
		return
	}
	respondJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

// Take in a writer, write out any response to json.
//
// example usecase -> frontend wants vehicles list, backend gets that from database,
// this function turns that database response to json, then that is sent to the frontend
func respondJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}

// format error responses as json to the frontend
func respondError(w http.ResponseWriter, status int, msg string) {
	respondJSON(w, status, map[string]string{"error": msg})
}
