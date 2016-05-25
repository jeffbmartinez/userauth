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

	cookie, err := getSIDCookie(r)
	if err != nil {
		log.WithError(err).Debug("Couldn't decode sid cookie")
		respond.Simple(w, http.StatusBadRequest)
		return
	}

	respond.JSON(w, cookie, http.StatusOK)
}

func getSIDCookie(r *http.Request) (model.SIDCookie, error) {
	var requestBody model.SessionInfoRequest
	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		log.WithError(err).Debug("Couldn't decode request body in SessionVerify")
		return model.SIDCookie{}, err
	}
	defer r.Body.Close()

	var cookie model.SIDCookie
	if err := safecookie.Get().Decode(model.SIDCookieName, requestBody.SID, &cookie); err != nil {
		log.WithError(err).WithFields(log.Fields{
			"cookie": requestBody.SID,
		}).Debug("Couldn't decode sid cookie")

		return model.SIDCookie{}, err
	}

	return cookie, nil
}
