package model

const (
	// SessionCookieName is the name of the cookie stored by the userauth service
	SessionCookieName = "session_info"
)

// SessionCookie contains the session information for a user's session.
type SessionCookie struct {
	SID string `json:"sid"` // Session ID
	UID string `json:"uid"` // User ID

	GoogleUserID string `json:"google_user_id"`
	HostedDomain string `json:"hosted_domain"`

	Email      string `json:"email"`
	Name       string `json:"name"`
	Picture    string `json:"picture"`
	GivenName  string `json:"given_name"`
	FamilyName string `json:"family_name"`
	Locale     string `json:"locale"`
}
