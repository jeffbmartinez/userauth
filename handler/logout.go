package handler

import (
	"net/http"

	"github.com/jeffbmartinez/respond"

	"github.com/jeffbmartinez/userauth/safecookie"
)

/*
Logout handles requests to log a user out. It works by expiring the session ID cookie, if
one exists.
*/
func Logout(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		respond.Simple(w, http.StatusMethodNotAllowed)
		return
	}

	// A cookie can't really be deleted. Overwriting an existing cookie of the same name
	// with a cookie that expires immediately is the next best thing.
	http.SetCookie(w, safecookie.GetExpiredSessionCookie())

	respond.Simple(w, http.StatusOK)
}
