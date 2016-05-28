package model

// SessionInfoRequest represents a request for /session/info
type SessionInfoRequest struct {
	SessionInfo string `json:"sessionInfo"` // encrypted session info cookie value (contains both sid and uid)
}
