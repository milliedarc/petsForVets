package httpapi

import (
	"net/http"

	"golang.org/x/crypto/bcrypt"

	"petsforvets/api/internal/db"
)

type appModeRequest struct {
	AppMode string `json:"app_mode"`
}

// PUT /api/auth/me/app-mode
//
// Switches between the pet-owner and clinic-team views. Only vets may use the
// clinic-team view; the frontend hides the option for everyone else, but the
// real check is here, because anyone can call the API directly.
func (s *Server) handleUpdateAppMode(w http.ResponseWriter, r *http.Request, user User) {
	var req appModeRequest
	if err := readJSON(w, r, &req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	switch req.AppMode {
	case "petOwner":
		// Everyone can be a pet owner.
	case "clinicTeam":
		if !user.Veterinary {
			writeError(w, http.StatusForbidden, "only clinic staff can use the clinic team view")
			return
		}
	default:
		writeError(w, http.StatusBadRequest, `app_mode must be "petOwner" or "clinicTeam"`)
		return
	}

	row, err := s.queries.UpdateUserAppMode(r.Context(), db.UpdateUserAppModeParams{
		ID:      user.ID,
		AppMode: req.AppMode,
	})
	if err != nil {
		serverError(w, r, err)
		return
	}
	writeJSON(w, http.StatusOK, User(row))
}

type changePasswordRequest struct {
	CurrentPassword string `json:"current_password"`
	NewPassword     string `json:"new_password"`
}

// POST /api/auth/change-password
//
// Checks the current password, saves the new one, and logs the user out of
// every device by deleting all their sessions. If someone had stolen a
// session cookie, changing the password now locks them out too.
func (s *Server) handleChangePassword(w http.ResponseWriter, r *http.Request, user User) {
	var req changePasswordRequest
	if err := readJSON(w, r, &req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if problem := checkPasswordRules(req.NewPassword); problem != "" {
		writeError(w, http.StatusBadRequest, problem)
		return
	}

	currentHash, err := s.queries.GetUserPasswordHash(r.Context(), user.ID)
	if err != nil {
		serverError(w, r, err)
		return
	}
	if bcrypt.CompareHashAndPassword([]byte(currentHash), []byte(req.CurrentPassword)) != nil {
		writeError(w, http.StatusBadRequest, "current password is incorrect")
		return
	}

	newHash, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		serverError(w, r, err)
		return
	}

	// Both statements succeed together or not at all. Without the transaction,
	// a crash between them could leave the new password set while old
	// sessions (possibly an attacker's) still work.
	err = s.inTx(r.Context(), func(q *db.Queries) error {
		if err := q.UpdateUserPassword(r.Context(), db.UpdateUserPasswordParams{
			ID:           user.ID,
			PasswordHash: string(newHash),
		}); err != nil {
			return err
		}
		return q.DeleteUserSessions(r.Context(), user.ID)
	})
	if err != nil {
		serverError(w, r, err)
		return
	}

	s.clearSessionCookie(w)
	w.WriteHeader(http.StatusNoContent)
}
