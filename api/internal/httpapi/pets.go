package httpapi

import (
	"errors"
	"net/http"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"

	"petsforvets/api/internal/db"
)

const (
	dateFormat = "2006-01-02" // Go writes date layouts using this exact reference date

	// Postgres error codes we turn into friendly 400 responses.
	pgForeignKeyViolation = "23503"
	pgCheckViolation      = "23514"
)

// petResponse is the JSON shape of a pet. It's built from the database row
// rather than sending the row directly, so the API shape and the table
// layout can change independently (and owner_id stays private).
type petResponse struct {
	ID                uuid.UUID  `json:"id"`
	Name              string     `json:"name"`
	DateOfBirth       *string    `json:"date_of_birth"` // "YYYY-MM-DD" or null
	Gender            string     `json:"gender"`
	Neutered          bool       `json:"neutered"`
	Colour            string     `json:"colour"`
	Imported          bool       `json:"imported"`
	ImportCountryCode string     `json:"import_country_code"`
	MicrochipNumber   string     `json:"microchip_number"`
	AvatarURL         *string    `json:"avatar_url"` // null until avatar uploads are ported
	Species           speciesRef `json:"species"`
	Breed             *breedRef  `json:"breed"` // null if no breed
}

type speciesRef struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	PetTypeID string    `json:"pet_type_id"`
}

type breedRef struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

func newPetResponse(row db.PetDetail) petResponse {
	pet := petResponse{
		ID:                row.ID,
		Name:              row.Name,
		Gender:            row.Gender,
		Neutered:          row.Neutered,
		Colour:            row.Colour,
		Imported:          row.Imported,
		ImportCountryCode: row.ImportCountryCode,
		MicrochipNumber:   row.MicrochipNumber,
		Species: speciesRef{
			ID:        row.SpeciesID,
			Name:      row.SpeciesName,
			PetTypeID: row.PetTypeID,
		},
	}
	if row.DateOfBirth != nil {
		dob := row.DateOfBirth.Format(dateFormat)
		pet.DateOfBirth = &dob
	}
	if row.BreedID != nil && row.BreedName != nil {
		pet.Breed = &breedRef{ID: *row.BreedID, Name: *row.BreedName}
	}
	return pet
}

// petRequest is the JSON body for creating or updating a pet.
// Note what's missing: owner_id. The owner is always the logged-in user,
// taken from the session. Never trust the client to say who it is.
type petRequest struct {
	Name              string     `json:"name"`
	SpeciesID         uuid.UUID  `json:"species_id"`
	BreedID           *uuid.UUID `json:"breed_id"`
	DateOfBirth       *string    `json:"date_of_birth"`
	Gender            string     `json:"gender"`
	Neutered          bool       `json:"neutered"`
	Colour            string     `json:"colour"`
	Imported          bool       `json:"imported"`
	ImportCountryCode string     `json:"import_country_code"`
	MicrochipNumber   string     `json:"microchip_number"`
}

// petFields is a petRequest that has passed validation.
type petFields struct {
	name, gender, colour, importCountryCode, microchipNumber string
	speciesID                                                uuid.UUID
	breedID                                                  *uuid.UUID
	dateOfBirth                                              *time.Time
	neutered, imported                                       bool
}

// validate checks the request and returns cleaned-up values, or a message
// explaining the first problem. The database also enforces most of these
// rules (NOT NULL, CHECK, foreign keys): this layer exists to give clear
// error messages instead of generic failures.
func (req petRequest) validate() (petFields, string) {
	f := petFields{
		name:              strings.TrimSpace(req.Name),
		speciesID:         req.SpeciesID,
		breedID:           req.BreedID,
		gender:            req.Gender,
		neutered:          req.Neutered,
		colour:            strings.TrimSpace(req.Colour),
		imported:          req.Imported,
		importCountryCode: strings.ToUpper(strings.TrimSpace(req.ImportCountryCode)),
		microchipNumber:   strings.TrimSpace(req.MicrochipNumber),
	}

	switch {
	case f.name == "":
		return f, "name is required"
	case utf8.RuneCountInString(f.name) > 100:
		return f, "name must be at most 100 characters"
	case f.speciesID == uuid.Nil:
		return f, "species_id is required"
	case f.gender != "" && f.gender != "male" && f.gender != "female":
		return f, `gender must be "male", "female" or empty`
	case utf8.RuneCountInString(f.colour) > 100:
		return f, "colour must be at most 100 characters"
	case len(f.microchipNumber) > 50:
		return f, "microchip_number must be at most 50 characters"
	case len(f.importCountryCode) > 3:
		return f, "import_country_code must be a country code like GB"
	}

	if !f.imported {
		f.importCountryCode = "" // a country only makes sense for imported pets
	}

	if req.DateOfBirth != nil && *req.DateOfBirth != "" {
		dob, err := time.Parse(dateFormat, *req.DateOfBirth)
		if err != nil {
			return f, "date_of_birth must look like 2020-01-31"
		}
		if dob.After(time.Now()) {
			return f, "date_of_birth can't be in the future"
		}
		f.dateOfBirth = &dob
	}

	return f, ""
}

// readPetRequest decodes and validates the body, writing a 400 if it's bad.
// It returns ok=false when the caller should stop.
func readPetRequest(w http.ResponseWriter, r *http.Request) (petFields, bool) {
	var req petRequest
	if err := readJSON(w, r, &req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return petFields{}, false
	}
	fields, problem := req.validate()
	if problem != "" {
		writeError(w, http.StatusBadRequest, problem)
		return petFields{}, false
	}
	return fields, true
}

// writePetDBError turns constraint violations into 400s and anything else into a 500.
func writePetDBError(w http.ResponseWriter, r *http.Request, err error) {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		switch pgErr.Code {
		case pgForeignKeyViolation:
			writeError(w, http.StatusBadRequest, "unknown species, or the breed doesn't belong to that species")
			return
		case pgCheckViolation:
			writeError(w, http.StatusBadRequest, "invalid pet details")
			return
		}
	}
	serverError(w, r, err)
}

// uuidFromPath reads {id} from the URL. A malformed id can't match anything,
// so it gets the same 404 as any other unknown id.
func uuidFromPath(w http.ResponseWriter, r *http.Request, notFoundMessage string) (uuid.UUID, bool) {
	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		writeError(w, http.StatusNotFound, notFoundMessage)
		return uuid.Nil, false
	}
	return id, true
}

func petIDFromPath(w http.ResponseWriter, r *http.Request) (uuid.UUID, bool) {
	return uuidFromPath(w, r, "pet not found")
}

// GET /api/pets
func (s *Server) handleListPets(w http.ResponseWriter, r *http.Request, user User) {
	rows, err := s.queries.ListPetsByOwner(r.Context(), user.ID)
	if err != nil {
		serverError(w, r, err)
		return
	}

	// make(..., 0, n) so an owner with no pets gets [] in the JSON, not null.
	pets := make([]petResponse, 0, len(rows))
	for _, row := range rows {
		pets = append(pets, newPetResponse(row))
	}
	writeJSON(w, http.StatusOK, pets)
}

// GET /api/pets/{id}
func (s *Server) handleGetPet(w http.ResponseWriter, r *http.Request, user User) {
	id, ok := petIDFromPath(w, r)
	if !ok {
		return
	}
	s.writePet(w, r, http.StatusOK, id, user.ID)
}

// POST /api/pets
func (s *Server) handleCreatePet(w http.ResponseWriter, r *http.Request, user User) {
	f, ok := readPetRequest(w, r)
	if !ok {
		return
	}

	id, err := s.queries.CreatePet(r.Context(), db.CreatePetParams{
		OwnerID:           user.ID,
		Name:              f.name,
		SpeciesID:         f.speciesID,
		BreedID:           f.breedID,
		DateOfBirth:       f.dateOfBirth,
		Gender:            f.gender,
		Neutered:          f.neutered,
		Colour:            f.colour,
		Imported:          f.imported,
		ImportCountryCode: f.importCountryCode,
		MicrochipNumber:   f.microchipNumber,
	})
	if err != nil {
		writePetDBError(w, r, err)
		return
	}
	s.writePet(w, r, http.StatusCreated, id, user.ID)
}

// PUT /api/pets/{id}
func (s *Server) handleUpdatePet(w http.ResponseWriter, r *http.Request, user User) {
	id, ok := petIDFromPath(w, r)
	if !ok {
		return
	}
	f, ok := readPetRequest(w, r)
	if !ok {
		return
	}

	updated, err := s.queries.UpdatePet(r.Context(), db.UpdatePetParams{
		ID:                id,
		OwnerID:           user.ID,
		Name:              f.name,
		SpeciesID:         f.speciesID,
		BreedID:           f.breedID,
		DateOfBirth:       f.dateOfBirth,
		Gender:            f.gender,
		Neutered:          f.neutered,
		Colour:            f.colour,
		Imported:          f.imported,
		ImportCountryCode: f.importCountryCode,
		MicrochipNumber:   f.microchipNumber,
	})
	if err != nil {
		writePetDBError(w, r, err)
		return
	}
	if updated == 0 {
		writeError(w, http.StatusNotFound, "pet not found")
		return
	}
	s.writePet(w, r, http.StatusOK, id, user.ID)
}

// DELETE /api/pets/{id}
func (s *Server) handleDeletePet(w http.ResponseWriter, r *http.Request, user User) {
	id, ok := petIDFromPath(w, r)
	if !ok {
		return
	}

	deleted, err := s.queries.DeletePet(r.Context(), db.DeletePetParams{ID: id, OwnerID: user.ID})
	if err != nil {
		serverError(w, r, err)
		return
	}
	if deleted == 0 {
		writeError(w, http.StatusNotFound, "pet not found")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// writePet loads one pet (with species and breed) and sends it as JSON.
func (s *Server) writePet(w http.ResponseWriter, r *http.Request, status int, id, ownerID uuid.UUID) {
	row, err := s.queries.GetPetForOwner(r.Context(), db.GetPetForOwnerParams{ID: id, OwnerID: ownerID})
	if errors.Is(err, pgx.ErrNoRows) {
		// Also what another user's pet looks like: we don't reveal it exists.
		writeError(w, http.StatusNotFound, "pet not found")
		return
	}
	if err != nil {
		serverError(w, r, err)
		return
	}
	writeJSON(w, status, newPetResponse(row))
}
