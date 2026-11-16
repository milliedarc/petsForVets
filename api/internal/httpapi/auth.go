package httpapi

import (
	"crypto/rand"
	"crypto/sha256"
	"errors"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"golang.org/x/crypto/bcrypt"

	"petsforvets/api/internal/db"
)

const (
	sessionCookieName = "session"
	sessionDuration   = 30 * 24 * time.Hour

	// Postgres error code for "duplicate value in a UNIQUE column".
	pgUniqueViolation = "23505"
)

// User is what the API sends to the frontend. It deliberately has no
// password hash field, so one can never be sent by accident.
// The `json:"..."` tags set the key names in the JSON output.
type User struct {
	ID         uuid.UUID `json:"id"`
	Email      string    `json:"email"`
	Name       string    `json:"name"`
	Surname    string    `json:"surname"`
	AppMode    string    `json:"app_mode"`
	Veterinary bool      `json:"veterinary"`
}

// A precomputed hash used when the email doesn't exist (see handleLogin).
var dummyHash, _ = bcrypt.GenerateFromPassword([]byte("not-a-real-password"), bcrypt.DefaultCost)

type registerRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Name     string `json:"name"`
	Surname  string `json:"surname"`
}

func (s *Server) handleRegister(w http.ResponseWriter, r *http.Request) {
	var req registerRequest
	if err := readJSON(w, r, &req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	email := normaliseEmail(req.Email)
	if !strings.Contains(email, "@") {
		writeError(w, http.StatusBadRequest, "a valid email is required")
		return
	}
	if problem := checkPasswordRules(req.Password); problem != "" {
		writeError(w, http.StatusBadRequest, problem)
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		serverError(w, r, err)
		return
	}

	row, err := s.queries.CreateUser(r.Context(), db.CreateUserParams{
		Email:        email,
		PasswordHash: string(hash),
		Name:         strings.TrimSpace(req.Name),
		Surname:      strings.TrimSpace(req.Surname),
	})
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgUniqueViolation {
			writeError(w, http.StatusConflict, "an account with this email already exists")
			return
		}
		serverError(w, r, err)
		return
	}

	if err := s.startSession(w, r, row.ID); err != nil {
		serverError(w, r, err)
		return
	}
	// CreateUserRow has exactly the same fields as User, so Go lets us convert.
	writeJSON(w, http.StatusCreated, User(row))
}

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (s *Server) handleLogin(w http.ResponseWriter, r *http.Request) {
	var req loginRequest
	if err := readJSON(w, r, &req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	user, err := s.queries.GetUserByEmail(r.Context(), normaliseEmail(req.Email))
	if errors.Is(err, pgx.ErrNoRows) {
		// No such user. We still run bcrypt so this answer takes as long as a
		// wrong password; otherwise response times reveal which emails exist.
		_ = bcrypt.CompareHashAndPassword(dummyHash, []byte(req.Password))
		writeError(w, http.StatusUnauthorized, "incorrect email or password")
		return
	}
	if err != nil {
		serverError(w, r, err)
		return
	}

	if bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)) != nil {
		writeError(w, http.StatusUnauthorized, "incorrect email or password")
		return
	}

	if err := s.startSession(w, r, user.ID); err != nil {
		serverError(w, r, err)
		return
	}
	writeJSON(w, http.StatusOK, User{
		ID:         user.ID,
		Email:      user.Email,
		Name:       user.Name,
		Surname:    user.Surname,
		AppMode:    user.AppMode,
		Veterinary: user.Veterinary,
	})
}

func (s *Server) handleLogout(w http.ResponseWriter, r *http.Request) {
	if cookie, err := r.Cookie(sessionCookieName); err == nil {
		if err := s.queries.DeleteSession(r.Context(), hashToken(cookie.Value)); err != nil {
			serverError(w, r, err)
			return
		}
	}
	s.clearSessionCookie(w)
	w.WriteHeader(http.StatusNoContent)
}

// clearSessionCookie tells the browser to delete the session cookie.
func (s *Server) clearSessionCookie(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:     sessionCookieName,
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   s.secureCookie,
		SameSite: http.SameSiteLaxMode,
	})
}

func (s *Server) handleMe(w http.ResponseWriter, r *http.Request, user User) {
	writeJSON(w, http.StatusOK, user)
}

// startSession creates a session row and sets the cookie. Like PHP's
// session_start(), but we own the storage: the sessions table.
func (s *Server) startSession(w http.ResponseWriter, r *http.Request, userID uuid.UUID) error {
	// Housekeeping: logging in is a good moment to clear out old sessions.
	if err := s.queries.DeleteExpiredSessions(r.Context()); err != nil {
		slog.Warn("deleting expired sessions", "err", err)
	}

	token := rand.Text() // 128 random bits, e.g. "4ZQ3KX7M2N..."
	expires := time.Now().Add(sessionDuration)

	err := s.queries.CreateSession(r.Context(), db.CreateSessionParams{
		TokenHash: hashToken(token),
		UserID:    userID,
		ExpiresAt: expires,
	})
	if err != nil {
		return err
	}

	http.SetCookie(w, &http.Cookie{
		Name:     sessionCookieName,
		Value:    token,
		Path:     "/",
		Expires:  expires,
		HttpOnly: true,                 // JavaScript can't read it, so XSS can't steal it
		Secure:   s.secureCookie,       // only sent over HTTPS in production
		SameSite: http.SameSiteLaxMode, // not sent on cross-site POSTs
	})
	return nil
}

// authedHandler is a handler that needs a logged-in user.
type authedHandler func(w http.ResponseWriter, r *http.Request, user User)

// requireUser checks the session cookie before calling next, and passes it
// the logged-in user. Roughly the PHP habit of
//
//	if (!isset($_SESSION['user_id'])) { http_response_code(401); exit; }
//
// at the top of each page, written once.
func (s *Server) requireUser(next authedHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie(sessionCookieName)
		if err != nil {
			writeError(w, http.StatusUnauthorized, "not logged in")
			return
		}

		row, err := s.queries.GetSessionUser(r.Context(), hashToken(cookie.Value))
		if errors.Is(err, pgx.ErrNoRows) {
			writeError(w, http.StatusUnauthorized, "session expired, please log in again")
			return
		}
		if err != nil {
			serverError(w, r, err)
			return
		}

		next(w, r, User(row))
	}
}

// checkPasswordRules returns a message if the password isn't acceptable.
// Shared by registration and password changes so the rules can't drift apart.
func checkPasswordRules(password string) string {
	switch {
	case len(password) < 8:
		return "password must be at least 8 characters"
	case len(password) > 72: // bcrypt only looks at the first 72 bytes
		return "password must be at most 72 characters"
	}
	return ""
}

func hashToken(token string) []byte {
	sum := sha256.Sum256([]byte(token))
	return sum[:]
}

func normaliseEmail(email string) string {
	return strings.ToLower(strings.TrimSpace(email))
}
