package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5"

	"github.com/mingassing/app/backend/internal/domain"
)

// Function to create a vehicle given an http request
//
// Likely will only be callable by developers, because we will need things
// like drag coeffs (i assume), maybe theres a pipeline we can write where
// users submit 3d scans of their cars using their phone cams, and we run phys
// sims on it to get the drag + other needed vars
func (h *handler) createVehicle(w http.ResponseWriter, r *http.Request) {
	var v domain.Vehicle
	// decode the json given by the http request
	if err := json.NewDecoder(r.Body).Decode(&v); err != nil {
		respondError(w, http.StatusBadRequest, "invalid json")
		return
	}
	// ensure the given values are valid
	if v.Name == "" || v.MassKg <= 0 || v.DragCoefficient <= 0 || v.FrontalAreaM2 <= 0 {
		respondError(w, http.StatusBadRequest, "name, mass_kg, drag_coefficient and frontal_area_m2 are required and must be positive")
		return
	}
	// create the vehicle in the database
	if err := h.db.CreateVehicle(r.Context(), &v); err != nil {
		respondError(w, http.StatusInternalServerError, "failed to create vehicle")
		return
	}
	// respond with the success status, + the vehicle data
	respondJSON(w, http.StatusCreated, v)
}

// Function to list all vehicles in the database
//
// used by user interface (User facing!) to list cars when user is selecting their model,
// also I imagine it will be used by product website (Dual use!) to get .length and see
// how many car models we support (and maybe even a "Do we support your car?" section
// of the marketing website)
func (h *handler) listVehicles(w http.ResponseWriter, r *http.Request) {
	vehicles, err := h.db.ListVehicles(r.Context())
	if err != nil {
		respondError(w, http.StatusInternalServerError, "failed to list vehicles")
		return
	}
	respondJSON(w, http.StatusOK, vehicles)
}

// Function to get a vehicle from the db from a specific id
//
// Used to get specific information of a vehicle without returning all vehicles
// to the frontend. For user profiles with saved vehicles under their account
func (h *handler) getVehicle(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid id")
		return
	}
	v, err := h.db.GetVehicle(r.Context(), id)
	if errors.Is(err, pgx.ErrNoRows) {
		respondError(w, http.StatusNotFound, "vehicle not found")
		return
	}
	if err != nil {
		respondError(w, http.StatusInternalServerError, "failed to get vehicle")
		return
	}
	respondJSON(w, http.StatusOK, v)
}

// Function to remove vehicles from the database.
//
// I doubt this will even be used, let alone used outside of an admin / developer
// sense, but its nice to have for CRUD (create, read, update, delete)
func (h *handler) deleteVehicle(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid id")
		return
	}
	deleted, err := h.db.DeleteVehicle(r.Context(), id)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "failed to delete vehicle")
		return
	}
	if !deleted {
		respondError(w, http.StatusNotFound, "vehicle not found")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
