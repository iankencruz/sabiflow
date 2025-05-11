package sessions

import (
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/sessions"
)

// store is the global cookie session store.
// It uses a HMAC secret to sign session cookies and prevent tampering.
var store *sessions.CookieStore

func init() {
	secret := os.Getenv("SESSION_KEY")
	if secret == "" {
		panic("SESSION_KEY environment variable is not set")
	}

	store = sessions.NewCookieStore([]byte(secret))
	store.Options = &sessions.Options{
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   86400 * 7, // 7 days
	}
}

// SetUserID stores the user ID in the session using the "user_session" cookie.
func SetUserID(w http.ResponseWriter, r *http.Request, userID int32) error {
	session, err := store.Get(r, "user_session")
	if err != nil {
		return err
	}
	session.Values["user_id"] = strconv.Itoa(int(userID))
	return session.Save(r, w)
}

// GetUserID retrieves the user ID from the "user_session" cookie.
// It returns 0 if the session is invalid or missing.
func GetUserID(r *http.Request) (int32, error) {
	session, err := store.Get(r, "user_session")
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

// ClearUserSession deletes the session cookie, typically used on logout.
func ClearUserSession(w http.ResponseWriter, r *http.Request) error {
	session, err := store.Get(r, "user_session")
	if err != nil {
		return err
	}
	session.Options.MaxAge = -1 // Marks cookie for deletion
	return session.Save(r, w)
}
