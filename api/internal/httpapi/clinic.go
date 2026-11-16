package httpapi

// Endpoints for clinic staff: the request queue, patients, and writing
// prescriptions. All of them are wrapped in requireStaff.

import (
	"errors"
	"net/http"
	"slices"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"

	"petsforvets/api/internal/db"
)

// requireStaff is requireUser plus "and must be clinic staff". It wraps the
// handler twice: first the session check, then the staff check.
func (s *Server) requireStaff(next authedHandler) http.HandlerFunc {
	return s.requireUser(func(w http.ResponseWriter, r *http.Request, user User) {
		if !user.Veterinary {
			writeError(w, http.StatusForbidden, "only clinic staff can do this")
			return
		}
		next(w, r, user)
	})
}

// ---------------------------------------------------------------- responses

type clinicSummaryResponse struct {
	PendingRequests  int64 `json:"pending_requests"`
	ApprovedRequests int64 `json:"approved_requests"`
	ReadyRequests    int64 `json:"ready_requests"`
	Patients         int64 `json:"patients"`
	Clients          int64 `json:"clients"`
}

type ownerRef struct {
	ID      uuid.UUID `json:"id"`
	Name    string    `json:"name"`
	Surname string    `json:"surname"`
	Email   string    `json:"email"`
}

type clinicRequestResponse struct {
	ID              uuid.UUID `json:"id"`
	Status          string    `json:"status"`
	Comments        string    `json:"comments"`
	RejectionReason string    `json:"rejection_reason"`
	CreatedAt       time.Time `json:"created_at"`
	Prescription    struct {
		ID           uuid.UUID `json:"id"`
		Dosage       string    `json:"dosage"`
		Quantity     int32     `json:"quantity"`
		MedicineName string    `json:"medicine_name"`
	} `json:"prescription"`
	Pet struct {
		ID          uuid.UUID `json:"id"`
		Name        string    `json:"name"`
		SpeciesName string    `json:"species_name"`
	} `json:"pet"`
	Owner ownerRef `json:"owner"`
}

func newClinicRequestResponse(row db.ClinicRequestDetail) clinicRequestResponse {
	resp := clinicRequestResponse{
		ID:              row.ID,
		Status:          row.Status,
		Comments:        row.Comments,
		RejectionReason: row.RejectionReason,
		CreatedAt:       row.CreatedAt,
		Owner:           ownerRef{ID: row.OwnerID, Name: row.OwnerName, Surname: row.OwnerSurname, Email: row.OwnerEmail},
	}
	resp.Prescription.ID = row.PrescriptionID
	resp.Prescription.Dosage = row.Dosage
	resp.Prescription.Quantity = row.Quantity
	resp.Prescription.MedicineName = row.MedicineName
	resp.Pet.ID = row.PetID
	resp.Pet.Name = row.PetName
	resp.Pet.SpeciesName = row.SpeciesName
	return resp
}

// clinicPetResponse is the owner's petResponse plus who owns it. Embedding
// petResponse (no field name) puts its fields at the top level of the JSON.
type clinicPetResponse struct {
	petResponse
	Owner        ownerRef `json:"owner"`
	OpenRequests int64    `json:"open_requests"`
}

func newClinicPetResponse(pet db.PetDetail, ownerName, ownerSurname, ownerEmail string, openRequests int64) clinicPetResponse {
	return clinicPetResponse{
		petResponse:  newPetResponse(pet),
		Owner:        ownerRef{ID: pet.OwnerID, Name: ownerName, Surname: ownerSurname, Email: ownerEmail},
		OpenRequests: openRequests,
	}
}

type clinicPrescriptionResponse struct {
	ID           uuid.UUID   `json:"id"`
	Dosage       string      `json:"dosage"`
	Quantity     int32       `json:"quantity"`
	CreatedAt    time.Time   `json:"created_at"`
	Medicine     medicineRef `json:"medicine"`
	PrescribedBy *string     `json:"prescribed_by"` // staff member's name, null if unknown
}

// ----------------------------------------------------------------- summary

// GET /api/clinic/summary
func (s *Server) handleClinicSummary(w http.ResponseWriter, r *http.Request, _ User) {
	row, err := s.queries.GetClinicSummary(r.Context())
	if err != nil {
		serverError(w, r, err)
		return
	}
	writeJSON(w, http.StatusOK, clinicSummaryResponse(row))
}

// ---------------------------------------------------------- request queue

var requestStatuses = []string{"Pending", "Approved", "Ready", "Rejected"}

// allowedFrom says which statuses a request may move to a new status from.
// Pending → Approved → Ready is the normal path; staff may also go straight
// from Pending to Ready, and may reject anything not yet finished.
var allowedFrom = map[string][]string{
	"Approved": {"Pending"},
	"Ready":    {"Pending", "Approved"},
	"Rejected": {"Pending", "Approved"},
}

// GET /api/clinic/prescription-requests?status=Pending&pet_id=...
// Both filters are optional.
func (s *Server) handleListClinicRequests(w http.ResponseWriter, r *http.Request, _ User) {
	var params db.ListClinicRequestsParams

	if status := r.URL.Query().Get("status"); status != "" {
		if !slices.Contains(requestStatuses, status) {
			writeError(w, http.StatusBadRequest, "status must be Pending, Approved, Ready or Rejected")
			return
		}
		params.Status = &status
	}
	if petID := r.URL.Query().Get("pet_id"); petID != "" {
		id, err := uuid.Parse(petID)
		if err != nil {
			writeError(w, http.StatusBadRequest, "pet_id must be a valid id")
			return
		}
		params.PetID = &id
	}

	rows, err := s.queries.ListClinicRequests(r.Context(), params)
	if err != nil {
		serverError(w, r, err)
		return
	}
	requests := make([]clinicRequestResponse, 0, len(rows))
	for _, row := range rows {
		requests = append(requests, newClinicRequestResponse(row))
	}
	writeJSON(w, http.StatusOK, requests)
}

type updateRequestStatusRequest struct {
	Status          string  `json:"status"`
	RejectionReason string  `json:"rejection_reason"`
	Comments        *string `json:"comments"` // optional note to the owner; null leaves it unchanged
}

// PATCH /api/clinic/prescription-requests/{id}
// Body: {"status": "Approved"} or {"status": "Rejected", "rejection_reason": "..."}
func (s *Server) handleUpdateClinicRequest(w http.ResponseWriter, r *http.Request, user User) {
	requestID, ok := uuidFromPath(w, r, "request not found")
	if !ok {
		return
	}

	var req updateRequestStatusRequest
	if err := readJSON(w, r, &req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	from, ok := allowedFrom[req.Status]
	if !ok {
		writeError(w, http.StatusBadRequest, "status must be Approved, Ready or Rejected")
		return
	}
	req.RejectionReason = strings.TrimSpace(req.RejectionReason)
	if req.Status == "Rejected" && req.RejectionReason == "" {
		writeError(w, http.StatusBadRequest, "please give the owner a reason for rejecting the request")
		return
	}
	if utf8.RuneCountInString(req.RejectionReason) > 500 {
		writeError(w, http.StatusBadRequest, "the reason must be at most 500 characters")
		return
	}
	if req.Comments != nil {
		trimmed := strings.TrimSpace(*req.Comments)
		if utf8.RuneCountInString(trimmed) > 500 {
			writeError(w, http.StatusBadRequest, "the note must be at most 500 characters")
			return
		}
		req.Comments = &trimmed
	}

	_, err := s.queries.UpdateClinicRequestStatus(r.Context(), db.UpdateClinicRequestStatusParams{
		NewStatus:       req.Status,
		RejectionReason: req.RejectionReason,
		Comments:        req.Comments,
		StaffID:         user.ID,
		ID:              requestID,
		FromStatuses:    from,
	})
	if errors.Is(err, pgx.ErrNoRows) {
		// Either no such request, or it isn't in a status this change is
		// allowed from (someone else got there first). Look it up to say which.
		s.explainRequestNotUpdated(w, r, requestID)
		return
	}
	if err != nil {
		serverError(w, r, err)
		return
	}

	row, err := s.queries.GetClinicRequest(r.Context(), requestID)
	if err != nil {
		serverError(w, r, err)
		return
	}
	writeJSON(w, http.StatusOK, newClinicRequestResponse(row))
}

func (s *Server) explainRequestNotUpdated(w http.ResponseWriter, r *http.Request, requestID uuid.UUID) {
	current, err := s.queries.GetClinicRequest(r.Context(), requestID)
	switch {
	case errors.Is(err, pgx.ErrNoRows):
		writeError(w, http.StatusNotFound, "request not found")
	case err != nil:
		serverError(w, r, err)
	default:
		writeError(w, http.StatusConflict, "this request has already been "+statusPhrases[current.Status])
	}
}

// statusPhrases finish the sentence "this request has already been ...".
var statusPhrases = map[string]string{
	"Pending":  "requested",
	"Approved": "approved",
	"Ready":    "marked as ready for collection",
	"Rejected": "rejected",
}

// ---------------------------------------------------------------- patients

// GET /api/clinic/pets?search=bella
// Matches pet name, owner name or email, or microchip number (any part of it).
func (s *Server) handleListClinicPets(w http.ResponseWriter, r *http.Request, _ User) {
	var search *string
	if term := strings.TrimSpace(r.URL.Query().Get("search")); term != "" {
		pattern := "%" + escapeLike(term) + "%"
		search = &pattern
	}

	rows, err := s.queries.ListClinicPets(r.Context(), search)
	if err != nil {
		serverError(w, r, err)
		return
	}
	pets := make([]clinicPetResponse, 0, len(rows))
	for _, row := range rows {
		pets = append(pets, newClinicPetResponse(row.PetDetail, row.OwnerName, row.OwnerSurname, row.OwnerEmail, row.OpenRequests))
	}
	writeJSON(w, http.StatusOK, pets)
}

// escapeLike makes % and _ in a search term match themselves rather than
// acting as LIKE wildcards. (Postgres's default escape character is \.)
func escapeLike(term string) string {
	return strings.NewReplacer(`\`, `\\`, `%`, `\%`, `_`, `\_`).Replace(term)
}

// GET /api/clinic/pets/{id}
func (s *Server) handleGetClinicPet(w http.ResponseWriter, r *http.Request, _ User) {
	petID, ok := petIDFromPath(w, r)
	if !ok {
		return
	}
	row, err := s.queries.GetClinicPet(r.Context(), petID)
	if errors.Is(err, pgx.ErrNoRows) {
		writeError(w, http.StatusNotFound, "pet not found")
		return
	}
	if err != nil {
		serverError(w, r, err)
		return
	}
	writeJSON(w, http.StatusOK, newClinicPetResponse(row.PetDetail, row.OwnerName, row.OwnerSurname, row.OwnerEmail, row.OpenRequests))
}

// GET /api/clinic/pets/{id}/prescriptions
func (s *Server) handleListClinicPetPrescriptions(w http.ResponseWriter, r *http.Request, _ User) {
	petID, ok := petIDFromPath(w, r)
	if !ok {
		return
	}
	rows, err := s.queries.ListClinicPetPrescriptions(r.Context(), petID)
	if err != nil {
		serverError(w, r, err)
		return
	}
	prescriptions := make([]clinicPrescriptionResponse, 0, len(rows))
	for _, row := range rows {
		prescriptions = append(prescriptions, newClinicPrescriptionResponse(row))
	}
	writeJSON(w, http.StatusOK, prescriptions)
}

func newClinicPrescriptionResponse(row db.ListClinicPetPrescriptionsRow) clinicPrescriptionResponse {
	resp := clinicPrescriptionResponse{
		ID:        row.ID,
		Dosage:    row.Dosage,
		Quantity:  row.Quantity,
		CreatedAt: row.CreatedAt,
		Medicine:  medicineRef{ID: row.MedicineID, Name: row.MedicineName, Description: row.MedicineDescription},
	}
	if row.PrescribedByName != nil {
		name := strings.TrimSpace(*row.PrescribedByName + " " + deref(row.PrescribedBySurname))
		resp.PrescribedBy = &name
	}
	return resp
}

func deref(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

type createPrescriptionRequest struct {
	MedicineID uuid.UUID `json:"medicine_id"`
	Dosage     string    `json:"dosage"`
	Quantity   int32     `json:"quantity"`
}

// POST /api/clinic/pets/{id}/prescriptions
func (s *Server) handleCreatePrescription(w http.ResponseWriter, r *http.Request, user User) {
	petID, ok := petIDFromPath(w, r)
	if !ok {
		return
	}

	var req createPrescriptionRequest
	if err := readJSON(w, r, &req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	req.Dosage = strings.TrimSpace(req.Dosage)
	switch {
	case req.MedicineID == uuid.Nil:
		writeError(w, http.StatusBadRequest, "please choose a medicine")
		return
	case req.Dosage == "":
		writeError(w, http.StatusBadRequest, "please give the dosage instructions")
		return
	case utf8.RuneCountInString(req.Dosage) > 200:
		writeError(w, http.StatusBadRequest, "the dosage must be at most 200 characters")
		return
	case req.Quantity < 1 || req.Quantity > 1000:
		writeError(w, http.StatusBadRequest, "the quantity must be between 1 and 1000")
		return
	}

	id, err := s.queries.CreatePrescription(r.Context(), db.CreatePrescriptionParams{
		PetID:        petID,
		MedicineID:   req.MedicineID,
		Dosage:       req.Dosage,
		Quantity:     req.Quantity,
		PrescribedBy: &user.ID,
	})
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.Code == pgForeignKeyViolation {
		// Either the pet or the medicine doesn't exist; the constraint name says which.
		if strings.Contains(pgErr.ConstraintName, "medicine") {
			writeError(w, http.StatusBadRequest, "unknown medicine")
		} else {
			writeError(w, http.StatusNotFound, "pet not found")
		}
		return
	}
	if err != nil {
		serverError(w, r, err)
		return
	}

	// Return it in the same shape as the list, so the page can just add it.
	rows, err := s.queries.ListClinicPetPrescriptions(r.Context(), petID)
	if err != nil {
		serverError(w, r, err)
		return
	}
	for _, row := range rows {
		if row.ID == id {
			writeJSON(w, http.StatusCreated, newClinicPrescriptionResponse(row))
			return
		}
	}
	serverError(w, r, errors.New("created prescription not found"))
}

// GET /api/clinic/medicines
func (s *Server) handleListMedicines(w http.ResponseWriter, r *http.Request, _ User) {
	rows, err := s.queries.ListMedicines(r.Context())
	if err != nil {
		serverError(w, r, err)
		return
	}
	medicines := make([]medicineRef, 0, len(rows))
	for _, row := range rows {
		medicines = append(medicines, medicineRef{ID: row.ID, Name: row.Name, Description: row.Description})
	}
	writeJSON(w, http.StatusOK, medicines)
}
