package model

const (
	// SessionCookieName is the name of the cookie stored by the userauth service
	SessionCookieName = "session_info"
)

// SessionCookie contains the session information for a user's session.
type SessionCookie struct {
	SID string `json:"sid"` // Session ID
	UID string `json:"uid"` // User ID
}
