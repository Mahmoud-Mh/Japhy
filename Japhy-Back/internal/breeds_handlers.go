package internal

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type BreedsHandler struct {
	service *BreedsService
}

func NewBreedsHandler(service *BreedsService) *BreedsHandler {
	return &BreedsHandler{service: service}
}

// GetBreeds handles GET /breeds with optional search parameters
func (h *BreedsHandler) GetBreeds(w http.ResponseWriter, r *http.Request) {
	params := BreedSearchParams{
		Species: r.URL.Query().Get("species"),
		PetSize: r.URL.Query().Get("pet_size"),
	}

	if minWeightStr := r.URL.Query().Get("min_weight"); minWeightStr != "" {
		if minWeight, err := strconv.Atoi(minWeightStr); err == nil {
			params.MinWeight = &minWeight
		}
	}

	if maxWeightStr := r.URL.Query().Get("max_weight"); maxWeightStr != "" {
		if maxWeight, err := strconv.Atoi(maxWeightStr); err == nil {
			params.MaxWeight = &maxWeight
		}
	}

	breeds, err := h.service.GetAll(params)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(breeds)
}

// GetBreed handles GET /breeds/{id}
func (h *BreedsHandler) GetBreed(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid breed ID", http.StatusBadRequest)
		return
	}

	breed, err := h.service.GetByID(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if breed == nil {
		http.Error(w, "Breed not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(breed)
}

// CreateBreed handles POST /breeds
func (h *BreedsHandler) CreateBreed(w http.ResponseWriter, r *http.Request) {
	var breed Breed
	if err := json.NewDecoder(r.Body).Decode(&breed); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if err := h.service.Create(&breed); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(breed)
}

// UpdateBreed handles PUT /breeds/{id}
func (h *BreedsHandler) UpdateBreed(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid breed ID", http.StatusBadRequest)
		return
	}

	var breed Breed
	if err := json.NewDecoder(r.Body).Decode(&breed); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if err := h.service.Update(id, &breed); err != nil {
		if err.Error() == "breed with id "+strconv.Itoa(id)+" not found" {
			http.Error(w, err.Error(), http.StatusNotFound)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	breed.ID = id
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(breed)
}

// DeleteBreed handles DELETE /breeds/{id}
func (h *BreedsHandler) DeleteBreed(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid breed ID", http.StatusBadRequest)
		return
	}

	if err := h.service.Delete(id); err != nil {
		if err.Error() == "breed with id "+strconv.Itoa(id)+" not found" {
			http.Error(w, err.Error(), http.StatusNotFound)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
} 