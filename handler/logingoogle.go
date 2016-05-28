package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	log "github.com/Sirupsen/logrus"
	"github.com/jeffbmartinez/respond"
	uuid "github.com/satori/go.uuid"
	"github.com/spf13/viper"

	"github.com/jeffbmartinez/userauth/model"
	"github.com/jeffbmartinez/userauth/safecookie"
)

const (
	googleTokenInfoEndpoint = "https://www.googleapis.com/oauth2/v3/tokeninfo?id_token="
)

/*
LoginGoogleRequest represents the request body for /login/google
*/
type LoginGoogleRequest struct {
	IDToken string `json:"idtoken"`
}

/*
LoginGoogleResponse represents the response body for /login/google
*/
type LoginGoogleResponse struct {
	Success bool `json:"success"`
}

/*
LoginGoogle allows a user to log in using their Google credentials. If the login succeeds
a session ID cookie is placed on the user's browser. As long as this cookie is present in
subsequent requests the user is considered logged in.
*/
func LoginGoogle(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		respond.Simple(w, http.StatusMethodNotAllowed)
		return
	}

	userauthGoogleClientID := viper.GetString("googleClientID")
	if userauthGoogleClientID == "" {
		log.Error("Couldn't load google client ID, all token verification attempts will fail")
		respond.Simple(w, http.StatusInternalServerError)
		return
	}

	var requestBody LoginGoogleRequest
	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		log.WithError(err).Warn("Couldn't decode request body in LoginGoogle")
		respond.Simple(w, http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	googleURI := fmt.Sprintf("%s%s", googleTokenInfoEndpoint, requestBody.IDToken)
	googleResponse, err := http.Get(googleURI)
	if err != nil {
		log.WithError(err).Error("Problem reaching google tokeninfo endpoint")
		respond.Simple(w, http.StatusInternalServerError)
		return
	}

	var decodedIDToken model.GoogleTokenInfo
	if err := json.NewDecoder(googleResponse.Body).Decode(&decodedIDToken); err != nil {
		log.WithError(err).Error("Couldn't decode google tokeninfo response")
		respond.Simple(w, http.StatusInternalServerError)
		return
	}
	defer googleResponse.Body.Close()

	// Optional config setting. If set, the token being verified must
	// have a matching id in the hd (hosted domain) claim
	userauthGoogleHostedDomainID := viper.GetString("googleHostedDomainID")

	// For Google Apps for Work user: check the 'hd' claim
	// If either userauth or google's token expects a hosted domain ID, it should be checked
	if userauthGoogleHostedDomainID != "" || decodedIDToken.HD != "" {
		if userauthGoogleHostedDomainID != decodedIDToken.HD {
			log.WithFields(log.Fields{
				"received": decodedIDToken.HD,
				"expected": userauthGoogleHostedDomainID,
			}).Warn("Google hosted domain ID mismatch")

			response := LoginGoogleResponse{Success: false}
			respond.JSON(w, response, http.StatusOK)
			return
		}
	}

	// TODO:
	// Make db entry connecting google user ID to userauth ID so in the future
	// facebook and other OAuth2 logins will work as well.
	// (for now just use the email as the userauth user ID)

	// Make sure the token contained the same client ID issued to this app
	if decodedIDToken.AUD != userauthGoogleClientID {
		log.WithFields(log.Fields{
			"received": decodedIDToken.AUD,
			"expected": userauthGoogleClientID,
		}).Warn("Google client ID mismatch")

		response := LoginGoogleResponse{Success: false}
		respond.JSON(w, response, http.StatusOK)
		return
	}

	cookie := model.SessionCookie{
		SID: uuid.NewV4().String(),
		UID: decodedIDToken.Email,
	}

	encryptedCookie, err := safecookie.MakeEncryptedSessionCookie(cookie)
	if err != nil {
		log.WithError(err).WithFields(log.Fields{
			"googleUserId": decodedIDToken.SUB,
			"email":        decodedIDToken.Email,
			"sid":          cookie.SID,
			"uid":          cookie.UID,
		}).Error("Couldn't encode cookie values, can't create session")

		respond.Simple(w, http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, encryptedCookie)

	response := LoginGoogleResponse{Success: true}
	respond.JSON(w, response, http.StatusOK)
}
