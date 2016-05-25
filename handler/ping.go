package handler

import (
	"net/http"

	log "github.com/Sirupsen/logrus"
	"github.com/jeffbmartinez/respond"
)

/*
Ping returns "pong". It serves as a health check.
*/
func Ping(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		respond.Simple(w, http.StatusMethodNotAllowed)
		return
	}

	log.WithFields(log.Fields{
		"remoteAddress": r.RemoteAddr,
	}).Debug("Received ping")

	respond.String(w, "pong", http.StatusOK)
}
