package safecookie

import (
	"net/http"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/gorilla/securecookie"
	"github.com/spf13/viper"

	"github.com/jeffbmartinez/userauth/model"
)

var (
	safeCookie *securecookie.SecureCookie

	/* ExpiredSessionCookie is an expired cookie. It is used to "remove" a
	cookie from a user's browser. Since a cookie can't really be deleted
	overwriting an existing cookie of the same name with a cookie that
	expires immediately is the next best thing.
	Also setting the cooking value to an empty string to be safe. */
	expiredSessionCookie = http.Cookie{
		Name:  model.SessionCookieName,
		Value: "",
		Path:  "/",
		// Secure:   true, // TODO: https://bitbucket.org/jeffbmartinez/doer/issues/39/enable-secure-cookies-in-userauth
		HttpOnly: true,
		MaxAge:   0,               // Get cookie to expire now
		Expires:  time.Unix(0, 0), // Set cookie to expire in the past
	}
)

func init() {
	viper.BindEnv("googleClientID", "USERAUTH_GOOGLE_CLIENT_ID")
	viper.BindEnv("secureCookieHashKey", "USERAUTH_SECURE_COOKIE_HASH_KEY")
	viper.BindEnv("secureCookieBlockKey", "USERAUTH_SECURE_COOKIE_BLOCK_KEY")

	hashKey := viper.GetString("secureCookieHashKey")
	if hashKey == "" {
		log.Fatal("Secure cookie hash key has not been configured. User authentication is not possible.")
	}

	blockKey := viper.GetString("secureCookieBlockKey")
	if blockKey == "" {
		log.Fatal("Secure cookie block key has not been configured. User authentication is unsafe.")
	}

	safeCookie = securecookie.New([]byte(hashKey), []byte(blockKey))
}

/*
Get returns an initialized instance of gorilla mux's securecookie.
This makes it easier to share the same secure cookie instance.
*/
func Get() *securecookie.SecureCookie {
	return safeCookie
}

// MakeEncryptedSessionCookie returns an encrypted cookie, ready to be stored on a user's browser.
func MakeEncryptedSessionCookie(sessionCookie model.SessionCookie) (*http.Cookie, error) {
	encodedCookieValue, err := Get().Encode(model.SessionCookieName, sessionCookie)
	if err != nil {
		return nil, err
	}

	return &http.Cookie{
		Name:  model.SessionCookieName,
		Value: encodedCookieValue,
		Path:  "/",
		// Secure:   true, // TODO: https://bitbucket.org/jeffbmartinez/doer/issues/39/enable-secure-cookies-in-userauth
		HttpOnly: true,
	}, nil
}

// GetExpiredSessionCookie returns a cookie which can be set on the browser to
// effectively remove any existing session cookie. There is no reliable way to
// delete cookies from a browser, so overwriting with an empty cookie that
// expires in the past is the next best thing.
func GetExpiredSessionCookie() *http.Cookie {
	return &expiredSessionCookie
}
