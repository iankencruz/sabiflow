package sessions

import (
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/sessions"
)

type Manager struct {
	store *sessions.CookieStore
}

const sessionName = "user_session"

// NewManager creates and returns a new session manager.
func NewManager() *Manager {
	secret := os.Getenv("SESSION_KEY")
	if secret == "" {
		panic("SESSION_KEY environment variable is not set")
	}

	store := sessions.NewCookieStore([]byte(secret))
	store.Options = &sessions.Options{
		Path:     "/",
		HttpOnly: true,
		Secure:   false, // set to true in production
		SameSite: http.SameSiteLaxMode,
		MaxAge:   86400 * 7, // 7 days
	}

	return &Manager{store: store}
}

// SetUserID stores the user ID in the session.
func (m *Manager) SetUserID(w http.ResponseWriter, r *http.Request, userID int32) error {
	session, err := m.store.Get(r, sessionName)
	if err != nil {
		return err
	}
	session.Values["user_id"] = strconv.Itoa(int(userID))
	return session.Save(r, w)
}

// GetUserID retrieves the user ID from the session.
func (m *Manager) GetUserID(r *http.Request) (int32, error) {
	session, err := m.store.Get(r, sessionName)
	if err != nil {
		return 0, err
	}
	raw, ok := session.Values["user_id"].(string)
	if !ok {
		return 0, http.ErrNoCookie
	}
	id, err := strconv.Atoi(raw)
	return int32(id), err
}

// Clear deletes the session cookie.
func (m *Manager) Clear(w http.ResponseWriter, r *http.Request) error {
	session, err := m.store.Get(r, sessionName)
	if err != nil {
		return err
	}
	session.Options.MaxAge = -1
	return session.Save(r, w)
}
