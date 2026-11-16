package httpapi

import (
	"net/http"

	"github.com/google/uuid"
)

type speciesResponse struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	PetTypeID string    `json:"pet_type_id"`
}

type breedResponse struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	SpeciesID uuid.UUID `json:"species_id"`
	PetTypeID string    `json:"pet_type_id"`
	Generic   bool      `json:"generic"`
}

// GET /api/species
func (s *Server) handleListSpecies(w http.ResponseWriter, r *http.Request, _ User) {
	rows, err := s.queries.ListSpecies(r.Context())
	if err != nil {
		serverError(w, r, err)
		return
	}
	species := make([]speciesResponse, 0, len(rows))
	for _, row := range rows {
		species = append(species, speciesResponse(row))
	}
	writeJSON(w, http.StatusOK, species)
}

// GET /api/breeds
func (s *Server) handleListBreeds(w http.ResponseWriter, r *http.Request, _ User) {
	rows, err := s.queries.ListBreeds(r.Context())
	if err != nil {
		serverError(w, r, err)
		return
	}
	breeds := make([]breedResponse, 0, len(rows))
	for _, row := range rows {
		breeds = append(breeds, breedResponse(row))
	}
	writeJSON(w, http.StatusOK, breeds)
}
