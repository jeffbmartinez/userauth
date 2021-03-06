package model

/*
GoogleTokenInfo models the response from calling the google user sign-in token id.
For more info see
https://developers.google.com/identity/sign-in/web/backend-auth#verify-the-integrity-of-the-id-token
*/
type GoogleTokenInfo struct {
	ISS string `json:"iss"`          // Issuer
	SUB string `json:"sub"`          // Unique Google User ID
	AZP string `json:"azp"`          // ?
	AUD string `json:"aud"`          // Audience
	IAT string `json:"iat"`          // ?
	EXP string `json:"exp"`          // Expiration
	HD  string `json:"hd,omitempty"` // Hosted Domain

	Email         string `json:"email"`
	EmailVerified string `json:"email_verified"`

	Name       string `json:"name"`
	Picture    string `json:"picture"`
	GivenName  string `json:"given_name"`
	FamilyName string `json:"family_name"`
	Locale     string `json:"locale"`
}
