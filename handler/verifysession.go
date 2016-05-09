package handler

import (
	"encoding/json"
	"net/http"

	log "github.com/Sirupsen/logrus"
	"github.com/jeffbmartinez/respond"

	"github.com/jeffbmartinez/userauth/model"
	"github.com/jeffbmartinez/userauth/safecookie"
)

// VerifySessionRequest represents a request for /verify/session
type VerifySessionRequest struct {
	SID string `json:"sid"`
}

// VerifySessionResponse represents a response for /verify/session
type VerifySessionResponse struct {
	Valid  bool   `json:"valid"`
	Reason string `json:"reason"`
}

/*
VerifySession checks a session ID for validity. It's response can be used to determine
the validity of the supplied session ID.
*/
func VerifySession(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		respond.Simple(w, http.StatusMethodNotAllowed)
		return
	}

	var requestBody VerifySessionRequest
	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		log.WithError(err).Warn("Couldn't decode request body in VerifySession")
		respond.Simple(w, http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	var cookie model.SIDCookie
	if err := safecookie.Get().Decode(model.SIDCookieName, requestBody.SID, &cookie); err != nil {
		log.WithError(err).WithFields(log.Fields{
			"cookie": requestBody.SID,
		}).Warn("Couldn't decode sid cookie")

		respond.Simple(w, http.StatusBadRequest)
		return
	}

	log.WithFields(log.Fields{
		"cookie": cookie,
	}).Info("Decrypted cookie")

	response := VerifySessionResponse{
		Valid:  true,
		Reason: "",
	}
	respond.JSON(w, response, http.StatusNotImplemented)
}
