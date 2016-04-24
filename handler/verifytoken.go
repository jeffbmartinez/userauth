package handler

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/jeffbmartinez/respond"

	"github.com/jeffbmartinez/userauth/model"
)

const (
	userauthGoogleClientIDEnv       = "USERAUTH_GOOGLE_CLIENT_ID"
	userauthGoogleHostedDomainIDEnv = "USERAUTH_GOOGLE_HOSTED_DOMAIN_ID"

	googleTokenInfoEndpoint = "https://www.googleapis.com/oauth2/v3/tokeninfo?id_token="
)

type VerifyTokenRequest struct {
	IDToken string `json:"idtoken"`
}

type VerifyTokenResponse struct {
	Valid bool `json:"valid"`
}

func VerifyGoogleIDToken(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		respond.Simple(w, http.StatusMethodNotAllowed)
		return
	}

	userauthGoogleClientID := os.Getenv(userauthGoogleClientIDEnv)
	if userauthGoogleClientID == "" {
		message := fmt.Sprintf("Couldn't load environment variable (%v), all token verification attempts will fail", userauthGoogleClientIDEnv)
		log.Println(message)
		respond.Simple(w, http.StatusInternalServerError)
		return
	}

	// optional environment variable. If set, the token being verified must
	// have a matching id in the hd (hosted domain) claim
	userauthGoogleHostedDomainID := os.Getenv(userauthGoogleHostedDomainIDEnv)

	var requestBody VerifyTokenRequest
	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		log.Printf("Couldn't decode request body; error: %v\n", err)
		respond.Simple(w, http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	googleURI := fmt.Sprintf("%s%s", googleTokenInfoEndpoint, requestBody.IDToken)
	googleResponse, err := http.Get(googleURI)
	if err != nil {
		log.Printf("Problem reaching google tokeninfo endpoint (%v)\n", err)
		respond.Simple(w, http.StatusInternalServerError)
		return
	}

	var decodedIDToken model.GoogleTokenInfo
	if err := json.NewDecoder(googleResponse.Body).Decode(&decodedIDToken); err != nil {
		fmt.Printf("Couldn't decode google tokeninfo response (%v)\n", err)
		respond.Simple(w, http.StatusInternalServerError)
		return
	}
	defer googleResponse.Body.Close()

	response := VerifyTokenResponse{
		Valid: true,
	}

	// If the user is a Google Apps for Work user, also check the 'hd' claim
	if userauthGoogleHostedDomainID != "" && userauthGoogleHostedDomainID != decodedIDToken.HD {
		response.Valid = false

	}

	if decodedIDToken.AUD != userauthGoogleClientID {
		response.Valid = false
	}

	respond.JSON(w, response, http.StatusOK)
}
