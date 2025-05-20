package sessions

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Manager struct {
	DB *pgxpool.Pool
}

func NewManager(db *pgxpool.Pool) *Manager {
	return &Manager{DB: db}
}

const (
	sessionCookieName = "user_session"
	sessionLifespan   = 7 * 24 * time.Hour // 7 days
)

func (m *Manager) SetUserID(w http.ResponseWriter, r *http.Request, userID int32) error {
	sessionToken := generateSessionToken()
	expiry := time.Now().Add(sessionLifespan)
	userAgent := r.UserAgent()

	args := pgx.NamedArgs{
		"token":      sessionToken,
		"user_id":    userID,
		"expires_at": expiry,
		"user_agent": userAgent,
	}

	_, err := m.DB.Exec(r.Context(), `
		INSERT INTO auth.sessions (token, user_id, expires_at, user_agent)
		VALUES (@token, @user_id, @expires_at, @user_agent)
		ON CONFLICT (token) DO UPDATE SET 
			expires_at = EXCLUDED.expires_at,
			user_agent = EXCLUDED.user_agent
	`, args)

	if err != nil {
		fmt.Println("❌ failed to insert session:", err)
		return fmt.Errorf("failed to insert session: %w", err)
	}

	fmt.Println("⚠️ session args:", args)

	cookie := &http.Cookie{
		Name:     sessionCookieName,
		Value:    sessionToken,
		Path:     "/",
		HttpOnly: true,
		Secure:   false, // set true in production behind HTTPS
		SameSite: http.SameSiteLaxMode,
		Expires:  expiry,
	}
	http.SetCookie(w, cookie)
	return nil
}

func (m *Manager) GetUserID(r *http.Request) (int32, error) {
	cookie, err := r.Cookie(sessionCookieName)
	if err != nil {
		return 0, err
	}

	var userID int32
	var expiresAt time.Time

	args := pgx.NamedArgs{
		"token": cookie.Value,
		"now":   time.Now(),
	}

	err = m.DB.QueryRow(r.Context(), `
		SELECT user_id, expires_at FROM auth.sessions
		WHERE token = @token AND expires_at > @now
	`, args).Scan(&userID, &expiresAt)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return 0, nil // no session found
		}
		return 0, fmt.Errorf("query session failed: %w", err)
	}

	// Optionally: refresh session expiry on activity
	newExpiry := time.Now().Add(sessionLifespan)
	if newExpiry.Sub(expiresAt) > (10 * time.Minute) {
		m.refreshExpiry(r.Context(), cookie.Value, newExpiry)
	}

	return userID, nil
}

func (m *Manager) refreshExpiry(ctx context.Context, token string, newExpiry time.Time) {
	_, _ = m.DB.Exec(ctx, `
		UPDATE auth.sessions SET expires_at = @expires_at WHERE token = @token
	`, pgx.NamedArgs{
		"token":      token,
		"expires_at": newExpiry,
	})
}

func (m *Manager) Clear(w http.ResponseWriter, r *http.Request) error {
	cookie, err := r.Cookie(sessionCookieName)
	if err != nil {
		return nil // No session to clear
	}

	_, _ = m.DB.Exec(r.Context(), `
		DELETE FROM auth.sessions WHERE token = @token
	`, pgx.NamedArgs{"token": cookie.Value})

	// Expire the cookie
	http.SetCookie(w, &http.Cookie{
		Name:     sessionCookieName,
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		Expires:  time.Unix(0, 0),
		MaxAge:   -1,
	})
	return nil
}

func (m *Manager) Put(ctx context.Context, token string, key string, value string) error {
	_, err := m.DB.Exec(ctx, `
		INSERT INTO auth.session_data (token, key, value)
		VALUES (@token, @key, @value)
		ON CONFLICT (token, key) DO UPDATE SET value = EXCLUDED.value
	`, pgx.NamedArgs{
		"token": token,
		"key":   key,
		"value": value,
	})
	return err
}

func (m *Manager) GetString(ctx context.Context, token string, key string) (string, error) {
	var val string
	err := m.DB.QueryRow(ctx, `
		SELECT value FROM auth.session_data
		WHERE token = @token AND key = @key
	`, pgx.NamedArgs{
		"token": token,
		"key":   key,
	}).Scan(&val)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", nil
		}
		return "", err
	}
	return val, nil
}

// generateSessionToken generates a new session token
func generateSessionToken() string {
	return strconv.FormatInt(time.Now().UnixNano(), 36) // simplistic token for now
}
