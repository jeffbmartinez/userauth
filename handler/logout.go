package handler

import (
	"net/http"
	"time"

	"github.com/jeffbmartinez/respond"
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

	/* A cookie can't really be deleted. Overwriting an existing cookie of the same name
	with a cookie that expires immediately is the next best thing.
	Also set the value to an empty string to be safe. */
	expiredSidCookie := http.Cookie{
		Name:     "sid",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		MaxAge:   0,               // Get cookie to expire now
		Expires:  time.Unix(0, 0), // Set cookie to expire in the past
	}

	http.SetCookie(w, &expiredSidCookie)

	respond.Simple(w, http.StatusOK)
}
