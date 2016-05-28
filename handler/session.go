package handler

import (
	"encoding/json"
	"net/http"

	log "github.com/Sirupsen/logrus"
	"github.com/jeffbmartinez/respond"

	"github.com/jeffbmartinez/userauth/model"
	"github.com/jeffbmartinez/userauth/safecookie"
)

/*
SessionInfo returns the contents of a session cookie.
*/
func SessionInfo(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		respond.Simple(w, http.StatusMethodNotAllowed)
		return
	}

	cookie, err := getSessionCookie(r)
	if err != nil {
		log.WithError(err).Debug("Couldn't decode session info cookie")
		respond.Simple(w, http.StatusBadRequest)
		return
	}

	respond.JSON(w, cookie, http.StatusOK)
}

func getSessionCookie(r *http.Request) (model.SessionCookie, error) {
	var requestBody model.SessionInfoRequest
	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		log.WithError(err).Debug("Couldn't decode request body in SessionVerify")
		return model.SessionCookie{}, err
	}
	defer r.Body.Close()

	var cookie model.SessionCookie
	if err := safecookie.Get().Decode(model.SessionCookieName, requestBody.SessionInfo, &cookie); err != nil {
		log.WithError(err).WithFields(log.Fields{
			"cookie": requestBody.SessionInfo,
		}).Debug("Couldn't decode sid cookie")

		return model.SessionCookie{}, err
	}

	return cookie, nil
}
