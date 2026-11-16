package httpapi

import (
	"errors"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"

	"petsforvets/api/internal/db"
)

type prescriptionResponse struct {
	ID       uuid.UUID   `json:"id"`
	Dosage   string      `json:"dosage"`
	Quantity int32       `json:"quantity"`
	Medicine medicineRef `json:"medicine"`
}

type medicineRef struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
}

type prescriptionRequestResponse struct {
	ID              uuid.UUID `json:"id"`
	Status          string    `json:"status"`
	Comments        string    `json:"comments"`
	RejectionReason string    `json:"rejection_reason"`
	CreatedAt       time.Time `json:"created_at"` // JSON: "2026-11-16T10:41:09.123Z"
	PrescriptionID  uuid.UUID `json:"prescription_id"`
	MedicineName    string    `json:"medicine_name"`
	Dosage          string    `json:"dosage"`
}

func newPrescriptionRequestResponse(row db.PrescriptionRequestDetail) prescriptionRequestResponse {
	return prescriptionRequestResponse{
		ID:              row.ID,
		Status:          row.Status,
		Comments:        row.Comments,
		RejectionReason: row.RejectionReason,
		CreatedAt:       row.CreatedAt,
		PrescriptionID:  row.PrescriptionID,
		MedicineName:    row.MedicineName,
		Dosage:          row.Dosage,
	}
}

// GET /api/pets/{id}/prescriptions
// Another user's pet gives an empty list, which reveals nothing.
func (s *Server) handleListPrescriptions(w http.ResponseWriter, r *http.Request, user User) {
	petID, ok := uuidFromPath(w, r, "pet not found")
	if !ok {
		return
	}

	rows, err := s.queries.ListPrescriptionsForPet(r.Context(), db.ListPrescriptionsForPetParams{
		PetID:   petID,
		OwnerID: user.ID,
	})
	if err != nil {
		serverError(w, r, err)
		return
	}

	prescriptions := make([]prescriptionResponse, 0, len(rows))
	for _, row := range rows {
		prescriptions = append(prescriptions, prescriptionResponse{
			ID:       row.ID,
			Dosage:   row.Dosage,
			Quantity: row.Quantity,
			Medicine: medicineRef{
				ID:          row.MedicineID,
				Name:        row.MedicineName,
				Description: row.MedicineDescription,
			},
		})
	}
	writeJSON(w, http.StatusOK, prescriptions)
}

// GET /api/pets/{id}/prescription-requests
func (s *Server) handleListPrescriptionRequests(w http.ResponseWriter, r *http.Request, user User) {
	petID, ok := uuidFromPath(w, r, "pet not found")
	if !ok {
		return
	}

	rows, err := s.queries.ListPrescriptionRequestsForPet(r.Context(), db.ListPrescriptionRequestsForPetParams{
		PetID:   petID,
		OwnerID: user.ID,
	})
	if err != nil {
		serverError(w, r, err)
		return
	}

	requests := make([]prescriptionRequestResponse, 0, len(rows))
	for _, row := range rows {
		requests = append(requests, newPrescriptionRequestResponse(row))
	}
	writeJSON(w, http.StatusOK, requests)
}

// POST /api/prescriptions/{id}/requests
func (s *Server) handleCreatePrescriptionRequest(w http.ResponseWriter, r *http.Request, user User) {
	prescriptionID, ok := uuidFromPath(w, r, "prescription not found")
	if !ok {
		return
	}

	requestID, err := s.queries.CreatePrescriptionRequest(r.Context(), db.CreatePrescriptionRequestParams{
		PrescriptionID: prescriptionID,
		OwnerID:        user.ID,
	})
	if errors.Is(err, pgx.ErrNoRows) {
		// The INSERT ... SELECT found no prescription of yours with that id.
		writeError(w, http.StatusNotFound, "prescription not found")
		return
	}
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.Code == pgUniqueViolation {
		// prescription_requests_one_open_idx: there's already an open request.
		writeError(w, http.StatusConflict, "there is already an open request for this prescription")
		return
	}
	if err != nil {
		serverError(w, r, err)
		return
	}

	row, err := s.queries.GetPrescriptionRequestForOwner(r.Context(), db.GetPrescriptionRequestForOwnerParams{
		ID:      requestID,
		OwnerID: user.ID,
	})
	if err != nil {
		serverError(w, r, err)
		return
	}
	writeJSON(w, http.StatusCreated, newPrescriptionRequestResponse(row))
}

// DELETE /api/prescription-requests/{id}
func (s *Server) handleCancelPrescriptionRequest(w http.ResponseWriter, r *http.Request, user User) {
	requestID, ok := uuidFromPath(w, r, "request not found")
	if !ok {
		return
	}

	deleted, err := s.queries.CancelPrescriptionRequest(r.Context(), db.CancelPrescriptionRequestParams{
		ID:      requestID,
		OwnerID: user.ID,
	})
	if err != nil {
		serverError(w, r, err)
		return
	}
	if deleted == 1 {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	// Nothing deleted: either it isn't yours/doesn't exist, or it's no longer
	// pending. Look it up to give the right answer.
	_, err = s.queries.GetPrescriptionRequestForOwner(r.Context(), db.GetPrescriptionRequestForOwnerParams{
		ID:      requestID,
		OwnerID: user.ID,
	})
	switch {
	case errors.Is(err, pgx.ErrNoRows):
		writeError(w, http.StatusNotFound, "request not found")
	case err != nil:
		serverError(w, r, err)
	default:
		writeError(w, http.StatusConflict, "only pending requests can be cancelled")
	}
}
